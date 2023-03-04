package sameriver

import (
	"runtime"
	"sync"
)

// moves entities according to their velocity
type PhysicsSystem struct {
	w               *World
	physicsEntities *UpdatedEntityList
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (p *PhysicsSystem) GetComponentDeps() []string {
	// TODO: do something with mass
	// TODO: impart velocity to collided objects?
	return []string{"Vec2D,Position", "Vec2D,Velocity", "Vec2D,Box", "Float64,Mass"}
}

func (p *PhysicsSystem) LinkWorld(w *World) {
	p.w = w
	p.physicsEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"physical",
			w.em.components.BitArrayFromNames([]string{"Position", "Velocity", "Box", "Mass"})))
}

func (p *PhysicsSystem) Update(dt_ms float64) {
	p.ParallelUpdate(dt_ms)
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

				// the logic is simpler to read that way
				pos := e.GetVec2D("Position")
				box := e.GetVec2D("Box")
				pos.ShiftCenterToBottomLeft(*box)
				defer pos.ShiftBottomLeftToCenter(*box)

				// calculate velocity
				vel := e.GetVec2D("Velocity")
				dx := vel.X * dt_ms
				dy := vel.Y * dt_ms

				// motion in x
				if pos.X+dx < 0 {
					// max out on the left
					pos.X = 0
				} else if pos.X+box.X+dx > float64(p.w.Width) {
					// max out on the right
					pos.X = float64(p.w.Width) - box.X
				} else {
					// otherwise move in x freely
					pos.X += dx
				}

				// motion in y
				if pos.Y+dy < 0 {
					// max out on the bottom
					pos.Y = 0
				} else if pos.Y+box.Y+dy > float64(p.w.Height) {
					// max out on the top
					pos.Y = float64(p.w.Height) - box.Y
				} else {
					// otherwise move in y freely
					pos.Y += dy
				}
			}
			wg.Done()
		}(startIndex, endIndex)
	}

	wg.Wait()
}

func (p *PhysicsSystem) SingleThreadUpdate(dt_ms float64) {
	// note: there are no function calls in the below, so we won't
	// be preempted while computing physics (this is very good, get it over with)
	for _, e := range p.physicsEntities.entities {
		// the logic is simpler to read that way
		pos := e.GetVec2D("Position")
		box := e.GetVec2D("Box")
		pos.ShiftCenterToBottomLeft(*box)
		defer pos.ShiftBottomLeftToCenter(*box)
		// calculate velocity
		vel := e.GetVec2D("Velocity")
		dx := vel.X * dt_ms
		dy := vel.Y * dt_ms
		// motion in x
		if pos.X+dx < 0 {
			// max out on the left
			pos.X = 0
		} else if pos.X+box.X+dx > float64(p.w.Width) {
			// max out on the right
			pos.X = float64(p.w.Width) - box.X
		} else {
			// otherwise move in x freely
			pos.X += dx
		}
		// motion in y
		if pos.Y+dy < 0 {
			// max out on the bottom
			pos.Y = 0
		} else if pos.Y+box.Y+dy > float64(p.w.Height) {
			// max out on the top
			pos.Y = float64(p.w.Height) - box.Y
		} else {
			// otherwise move in y freely
			pos.Y += dy
		}
	}
}

func (p *PhysicsSystem) Expand(n int) {
	// nil?
}
