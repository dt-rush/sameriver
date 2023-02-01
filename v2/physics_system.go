// moves entities according to their velocity
package sameriver

import (
	"runtime"
	"sync"
)

type PhysicsSystem struct {
	w               *World
	physicsEntities *UpdatedEntityList
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (s *PhysicsSystem) GetComponentDeps() []string {
	return []string{"Vec2D,Position", "Vec2D,Velocity", "Vec2D,Box", "Float64,Mass"}
}

func (s *PhysicsSystem) LinkWorld(w *World) {
	s.w = w
	s.physicsEntities = w.em.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"physical",
			w.em.components.BitArrayFromNames([]string{"Position", "Velocity", "Box", "Mass"})))
}

func (s *PhysicsSystem) Update(dt_ms float64) {
	numWorkers := runtime.NumCPU()
	entitiesChunkSize := len(s.physicsEntities.entities) / numWorkers
	entitiesChunks := make([][]*Entity, numWorkers)
	for i := 0; i < numWorkers; i++ {
		start := i * entitiesChunkSize
		end := (i + 1) * entitiesChunkSize
		if i == numWorkers-1 {
			end = len(s.physicsEntities.entities)
		}
		entitiesChunks[i] = s.physicsEntities.entities[start:end]
	}
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(entities []*Entity) {
			defer wg.Done()
			for _, e := range entities {
				pos := e.GetVec2D("Position")
				box := e.GetVec2D("Box")
				// the logic is simpler to read with shifting to the corner
				// TODO: modify a local copy, don't actually modify the real box
				pos.ShiftCenterToBottomLeft(box)
				// unshift position back to center
				defer pos.ShiftBottomLeftToCenter(box)
				// calculate velocity
				vel := e.GetVec2D("Velocity")
				dx := vel.X * dt_ms
				dy := vel.Y * dt_ms
				// motion in x
				if pos.X+dx < 0 {
					// max out on the left
					pos.X = 0
				} else if pos.X+box.X+dx > float64(s.w.Width) {
					// max out on the right
					pos.X = float64(s.w.Width) - box.X
				} else {
					// otherwise move in x freely
					pos.X += dx
				}
				// motion in y
				if pos.Y+dy < 0 {
					// max out on the bottom
					pos.Y = 0
				} else if pos.Y+box.Y+dy > float64(s.w.Height) {
					// max out on the top
					pos.Y = float64(s.w.Height) - box.Y
				} else {
					// otherwise move in y freely
					pos.Y += dy
				}
			}
		}(entitiesChunks[i])
	}
	wg.Wait()
}
