// The collision detection in this sytem has 4 main parts:
//
// 1. a method to check collisions invoked by the game every game loop
//
// 2. an UpdatedEntityList of entities having Position and HitBox
//
// 3. a special data structure which holds rate limiters for each possible
// 	collision
//
//
// Datastructure (3.) - triangular rateLimiters array
//
// The rate limiters data structed is "collision-indexed", meaning it is indexed
// [i][j], where i and j are ID's and i < j. That is, each pairing of ID's
// is produced by matching each ID with all those greater than it.
//
// A collision-indexed data structure of ResettableRateLimiters
// used to avoid notifying of collisions too often. The need for this arises
// from the fact that we want to run the collision-checking logic as often as
// possible, but we don't want to send collision events at 30 times a second.
// These rate limiters rate-limit the sending of messages on a channel when we
// detect collisions, in order to save resources (internally they use a
// sync.Once which can be reset either by a natural delay or externally, in
// a goroutine-safe way)
//
//          j
//
//      0 1 2 3 4
//     0  r r r r
//     1    r r r
//  i  2      r r
//     3        r
//     4
//
package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type CollisionSystem struct {
	w                  *World
	collidableEntities *UpdatedEntityList
	rateLimiterArray   CollisionRateLimiterArray
	sh                 *SpatialHashSystem `sameriver-system-dependency:"-"`
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{rateLimiterArray: NewCollisionRateLimiterArray()}
}

func (s *CollisionSystem) LinkWorld(w *World) {
	s.w = w
	// query a regularly updated list of the entities which are collidable
	// (position and hitbox)
	s.collidableEntities = w.em.GetSortedUpdatedEntityList(
		EntityQueryFromComponentBitArray(
			"collidable",
			MakeComponentBitArray([]ComponentType{
				BOX_COMPONENT})))
	// add a callback to the UpdatedEntityList of collidable entities
	// so that whenever an entity is removed, we will reset its rate limiters
	// in the collision rate limiter array (to guard against an entity
	// despawning, a new entity spawning with its ID, and failing a collision
	// test (rare prehaps, but an edge case we nonetheless want to avoid)
	s.collidableEntities.AddCallback(
		func(signal EntitySignal) {
			if signal.SignalType == ENTITY_REMOVE {
				s.rateLimiterArray.Reset(signal.Entity)
			}
		})
}

// Iterates through the entities in the UpdatedEntityList using a handshake
// pattern, where, given a sorted list of ID's corresponding to collidable
// entities, i is compared with all ID's after i, then i + 1 is compared with
// all entities after i + 1, etc. (basically we iterate through the
// collision-indexed rate-limiter 2d triangular array row by row, left to right)
//
// If a collision is confirmed by checking their positions and bounding boxes,
// we attempt to send a collision event through the channel to be processed
// by goroutine 2 ("Event filtering and sending"), but we rate-limit sending
// events for each possible collision [i][j] using the rate limiter at [i][j]
// in rateLimiters, so if we already sent one within the timeout, we just move on.
func (s *CollisionSystem) Update(dt_ms float64) {

	entities := s.collidableEntities.entities

	// NOTE: The ID's in collidableEntities are in sorted order,
	// so the rateLimiterArray access condition that i < j is respected
	// check each possible collison between entities in the list by doing a
	// handshake pattern
	for ix := uint16(0); ix < uint16(len(entities)); ix++ {
		for jx := ix + 1; jx < uint16(len(entities)); jx++ {
			// get the entity ID's
			i := entities[ix]
			j := entities[jx]
			// check the collision
			if s.TestCollision(uint16(i.ID), uint16(j.ID)) {
				// if colliding, send the message (rate-limited)
				s.rateLimiterArray.GetRateLimiter(i.ID, j.ID).Do(
					func() {
						s.w.ev.Publish(COLLISION_EVENT,
							CollisionData{EntityA: i, EntityB: j})
					})
				// TODO: move both entities away from their common center?
				// generalized callback function probably best (with a set of
				// predefined ones)
			}
		}
	}
}

// Test collision between two entities
func (s *CollisionSystem) TestCollision(i uint16, j uint16) bool {
	iPos := s.w.em.Components.Position[i]
	iBox := s.w.em.Components.Box[i]
	jPos := s.w.em.Components.Position[j]
	jBox := s.w.em.Components.Box[j]
	iRect := sdl.Rect{
		int32(iPos.X),
		int32(iPos.Y),
		int32(iBox.X),
		int32(iBox.Y)}
	jRect := sdl.Rect{
		int32(jPos.X),
		int32(jPos.Y),
		int32(jBox.X),
		int32(jBox.Y)}
	return iRect.HasIntersection(&jRect)
}
