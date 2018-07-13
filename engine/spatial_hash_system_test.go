package engine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSpatialHashInsertion(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
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
			for _, entity := range (*table)[cell[0]][cell[1]] {
				if entity == e {
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

func TestSpatialHashLargeEntity(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	pos := Vec2D{0, 0}
	box := Vec2D{25, 25}
	cells := [][2]int{
		[2]int{0, 0},
		[2]int{1, 1},
		[2]int{1, 2},
		[2]int{2, 1},
		[2]int{2, 2},
	}
	e, _ := w.em.Spawn(spatialSpawnRequestData(pos, box))
	w.Update(FRAME_SLEEP_MS)
	table := sh.CurrentTablePointer()
	for _, cell := range cells {
		inCell := false
		for _, entity := range (*table)[cell[0]][cell[1]] {
			if entity == e {
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

// testing that the spatial hash will not double-compute while one calculation
// is already running
func TestSpatialHashDoubleCompute(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	// spawn a lot of entities
	for i := 0; i < MAX_ENTITIES; i++ {
		pos := Vec2D{100 * rand.Float64(), 100 * rand.Float64()}
		w.em.Spawn(spatialSpawnRequestData(Vec2D{1, 1}, pos))
	}
	N := 512
	// update the world numerous times at the same time
	for i := 0; i < N; i++ {
		go func() {
			sh.Update(FRAME_SLEEP_MS)
		}()
	}
	if sh.timesComputed.Load() == 512 {
		t.Fatal("spatial hash did not guard against starting compute while " +
			"already in progress")
	}
}

func TestSpatialHashTableCopy(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	w.em.Spawn(spatialSpawnRequestData(Vec2D{0, 0}, Vec2D{1, 1}))
	w.Update(FRAME_SLEEP_MS)
	table := sh.CurrentTablePointer()
	tableCopy := sh.CurrentTableCopy()
	if (*table)[0][0][0] != tableCopy[0][0][0] {
		t.Fatal("CurrentTableCopy() doesn't return a copy")
	}
}

func TestSpatialHashTableToString(t *testing.T) {
	w := NewWorld(1024, 1024)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	table := sh.CurrentTablePointer()
	table.String()
}
