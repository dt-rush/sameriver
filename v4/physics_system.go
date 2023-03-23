package sameriver

import (
	"runtime"
	"sync"
)

// moves entities according to their velocity
type PhysicsSystem struct {
	granularity     int
	w               *World
	physicsEntities *UpdatedEntityList
	h               *SpatialHasher
}

func NewPhysicsSystem() *PhysicsSystem {
	return NewPhysicsSystemWithGranularity(1)
}

func NewPhysicsSystemWithGranularity(granularity int) *PhysicsSystem {
	return &PhysicsSystem{
		granularity: granularity,
	}
}

func (p *PhysicsSystem) GetComponentDeps() []any {
	// TODO: do something with mass
	// TODO: impart momentum to collided objects?
	return []any{
		POSITION, VEC2D, "POSITION",
		VELOCITY, VEC2D, "VELOCITY",
		ACCELERATION, VEC2D, "ACCELERATION",
		BOX, VEC2D, "BOX",
		MASS, FLOAT64, "MASS",
	}
}

func (p *PhysicsSystem) LinkWorld(w *World) {
	p.w = w
	p.physicsEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"physical",
			w.em.components.BitArrayFromIDs([]ComponentID{POSITION, VELOCITY, ACCELERATION, BOX, MASS})))
	p.h = NewSpatialHasher(10, 10, w)
}

func (p *PhysicsSystem) Update(dt_ms float64) {
	p.h.Update()
	sum_dt := 0.0
	for i := 0; i < p.granularity; i++ {
		p.ParallelUpdate(dt_ms / float64(p.granularity))
		sum_dt += dt_ms / float64(p.granularity)
	}
}

func (p *PhysicsSystem) physics(e *Entity, dt_ms float64) {

	// the logic is simpler to read that way
	pos := e.GetVec2D(POSITION)
	box := e.GetVec2D(BOX)
	pos.ShiftCenterToBottomLeft(*box)
	defer pos.ShiftBottomLeftToCenter(*box)

	// calculate velocity
	acc := e.GetVec2D(ACCELERATION)
	vel := e.GetVec2D(VELOCITY)
	vel.X += acc.X * dt_ms
	vel.Y += acc.Y * dt_ms
	dx := vel.X * dt_ms
	dy := vel.Y * dt_ms

	// motion in x
	// max out on world border in x
	if pos.X+dx < 0 || pos.X+box.X+dx > float64(p.w.Width) {
		dx = 0
	} else {
		// otherwise move in x freely
		pos.X += dx
	}

	// motion in y
	// max out on world border in y
	if pos.Y+dy < 0 || pos.Y+box.Y+dy > float64(p.w.Height) {
		dy = 0
	} else {
		// otherwise move in y freely
		pos.Y += dy
	}

	// check collisions using spatial hasher
	// TODO: really we should check / resolve all collisions after applying dx,dy
	testCollision := func(i *Entity, j *Entity) bool {
		iPos := i.GetVec2D(POSITION)
		iBox := i.GetVec2D(BOX)
		jPos := j.GetVec2D(POSITION)
		jBox := j.GetVec2D(BOX)
		return RectIntersectsRect(*iPos, *iBox, *jPos, *jBox)
	}
	cellX0, cellX1, cellY0, cellY1 := p.h.CellRangeOfRect(*pos, *box)
	collided := false
	for y := cellY0; y <= cellY1 && !collided; y++ {
		for x := cellX0; x <= cellX1 && !collided; x++ {
			if x < 0 || x >= p.h.GridX || y < 0 || y >= p.h.GridY {
				continue
			}
			entities := p.h.Entities(x, y)
			for i := 0; i < len(entities) && !collided; i++ {
				other := entities[i]
				if other != e && testCollision(e, other) {
					// undo the action if a collision occurs
					pos.X -= dx
					pos.Y -= dy
					collided = true
					break
				}
			}
		}
	}
}

func (p *PhysicsSystem) ParallelUpdate(dt_ms float64) {
	// divide the entities into N segments,
	// where N is the number of CPU cores
	numWorkers := runtime.NumCPU()
	entitiesPerWorker := len(p.physicsEntities.entities) / numWorkers
	remainder := len(p.physicsEntities.entities) % numWorkers

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		startIndex := i * entitiesPerWorker
		endIndex := (i + 1) * entitiesPerWorker
		if i == numWorkers-1 {
			endIndex += remainder
		}

		go func(startIndex, endIndex int) {
			for j := startIndex; j < endIndex; j++ {
				e := p.physicsEntities.entities[j]
				p.physics(e, dt_ms)
			}
			wg.Done()
		}(startIndex, endIndex)
	}

	wg.Wait()
}

func (p *PhysicsSystem) SingleThreadUpdate(dt_ms float64) {
	// note: there are no function calls in the below, so we won't
	// be preempted while computing physics (this is very good, get it over with)
	for i := range p.physicsEntities.entities {
		e := p.physicsEntities.entities[i]
		p.physics(e, dt_ms)
	}
}

func (p *PhysicsSystem) Expand(n int) {
	// nil?
}
