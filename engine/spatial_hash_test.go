package engine

import (
	"fmt"
	"testing"
)

func TestSpatialHashInsertion(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystem(sh)
	testData := map[[2]Vec2D][][2]int{
		[2]Vec2D{Vec2D{0, 0}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{1, 1}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{4, 4}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{0, 11}, Vec2D{1, 1}}:  [][2]int{[2]int{0, 1}},
		[2]Vec2D{Vec2D{11, 11}, Vec2D{1, 1}}: [][2]int{[2]int{1, 1}},
		[2]Vec2D{Vec2D{41, 41}, Vec2D{1, 1}}: [][2]int{[2]int{4, 4}},
		[2]Vec2D{Vec2D{99, 99}, Vec2D{1, 1}}: [][2]int{[2]int{9, 9}},
		[2]Vec2D{Vec2D{11, 99}, Vec2D{1, 1}}: [][2]int{[2]int{1, 9}},
	}
	entityCells := make(map[*EntityToken][][2]int)
	for posbox, cells := range testData {
		e, _ := w.em.Spawn(spatialSpawnRequestData(posbox[0], posbox[1]))
		entityCells[e] = cells
	}
	w.Update(FRAME_SLEEP_MS)
	for e, cells := range entityCells {
		for _, cell := range cells {
			table := sh.CurrentTablePointer()
			inCell := false
			for _, entityGridPosition := range (*table)[cell[0]][cell[1]] {
				if e == entityGridPosition.entity {
					inCell = true
				}
			}
			if !inCell {
				t.Fatal(fmt.Sprintf("%v,%v was not mapped to cell %v",
					w.em.Components.Position[e.ID],
					w.em.Components.Box[e.ID],
					cell))
			}
		}
	}
}
