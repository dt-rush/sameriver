// The collision detection in this sytem has 4 main parts:
//
// 1. a method to check collisions invoked by the game every game loop
//
// 2. an UpdatedEntityList of entities having Position and HitBox
//
//  3. a special data structure which holds rate limiters for each possible
//     collision
//
// # Datastructure (3.) - triangular rateLimiters array
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
//	        j
//
//	    0 1 2 3 4
//	   0  r r r r
//	   1    r r r
//	i  2      r r
//	   3        r
//	   4
package sameriver

import (
	"time"
)

type CollisionData struct {
	This  *Entity
	Other *Entity
}

type CollisionSystem struct {
	w                  *World
	collidableEntities *UpdatedEntityList
	rateLimiterArray   CollisionRateLimiterArray
	sh                 *SpatialHashSystem `sameriver-system-dependency:"-"`
}

func NewCollisionSystem(rateLimit time.Duration) *CollisionSystem {
	return &CollisionSystem{
		rateLimiterArray: NewCollisionRateLimiterArray(rateLimit),
	}
}

func (s *CollisionSystem) GetComponentDeps() []string {
	return []string{"Vec2D,Position", "Vec2D,Box"}
}

func (s *CollisionSystem) LinkWorld(w *World) {
	s.w = w
	// Filter a regularly updated list of the entities which are collidable
	// (position and hitbox)
	s.collidableEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"collidable",
			w.em.components.BitArrayFromNames([]string{"Position", "Box"})))
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
	// NOTE: The ID's in collidableEntities are in sorted order,
	// so the rateLimiterArray access condition that i < j is respected
	// check each possible collison between entities in the list by doing a
	// handshake pattern
	for x := 0; x < s.sh.GridX; x++ {
		for y := 0; y < s.sh.GridY; y++ {
			s.checkEntities(s.sh.Table[x][y])
		}
	}
}

func (s *CollisionSystem) checkEntities(entities []*Entity) {
	// NOTE: we guard for despawns since the entities in the spatial hash
	// table might have been despawned since the last time a spatial hash
	// was computed (not every system is guaranteed to run every update loop,
	// so maybe spatial hash didn't run but an entity or world logic did, to
	// despawn one of the tokens still stored in the last-computed spatial hash
	// table).
	for ix := uint16(0); ix < uint16(len(entities)); ix++ {
		i := entities[ix]
		if i.Despawned {
			continue
		}
		for jx := ix + 1; jx < uint16(len(entities)); jx++ {
			j := entities[jx]
			if j.Despawned {
				continue
			}
			// required that i.ID , j.ID for the rate limiter array
			if j.ID < i.ID {
				j, i = i, j
			}
			if !s.rateLimiterArray.GetRateLimiter(i.ID, j.ID).Limited() &&
				s.TestCollision(i, j) {
				s.DoCollide(i, j)
			}
		}
	}
}

func (s *CollisionSystem) DoCollide(i *Entity, j *Entity) {
	s.rateLimiterArray.GetRateLimiter(i.ID, j.ID).Do(
		func() {
			s.w.Events.Publish("collision",
				CollisionData{This: i, Other: j})
			s.w.Events.Publish("collision",
				CollisionData{This: j, Other: i})
		})
}

// Test collision between two entities
func (s *CollisionSystem) TestCollision(i *Entity, j *Entity) bool {
	iPos := i.GetVec2D("Position")
	iBox := i.GetVec2D("Box")
	jPos := j.GetVec2D("Position")
	jBox := j.GetVec2D("Box")
	intersects := RectIntersectsRect(iPos, iBox, jPos, jBox)
	return intersects
}
