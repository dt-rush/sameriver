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
		[2]Vec2D{Vec2D{5, 5}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{1, 1}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{4, 4}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{1, 11}, Vec2D{1, 1}}:  [][2]int{[2]int{0, 1}},
		[2]Vec2D{Vec2D{11, 11}, Vec2D{1, 1}}: [][2]int{[2]int{1, 1}},
		[2]Vec2D{Vec2D{41, 41}, Vec2D{1, 1}}: [][2]int{[2]int{4, 4}},
		[2]Vec2D{Vec2D{99, 99}, Vec2D{1, 1}}: [][2]int{[2]int{9, 9}},
		[2]Vec2D{Vec2D{11, 99}, Vec2D{1, 1}}: [][2]int{[2]int{1, 9}},
	}
	entityCells := make(map[*Entity][][2]int)
	for posbox, cells := range testData {
		e, _ := testingSpawnSpatial(w, posbox[0], posbox[1])
		entityCells[e] = cells
	}
	w.Update(FRAME_SLEEP_MS / 2)
	for e, cells := range entityCells {
		for _, cell := range cells {
			inCell := false
			for _, entity := range sh.Table[cell[0]][cell[1]] {
				if entity == e {
					inCell = true
				}
			}
			if !inCell {
				t.Fatal(fmt.Sprintf("%v,%v was not mapped to cell %v",
					e.GetVec2D("Position"),
					e.GetVec2D("Box"),
					cell))
			}
		}
	}
}

func TestSpatialHashLargeEntity(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	pos := Vec2D{20, 20}
	box := Vec2D{5, 5}
	cells := [][2]int{
		[2]int{1, 1},
		[2]int{1, 2},
		[2]int{2, 1},
		[2]int{2, 2},
	}
	e, _ := testingSpawnSpatial(w, pos, box)
	w.Update(FRAME_SLEEP_MS / 2)
	for _, cell := range cells {
		inCell := false
		for _, entity := range sh.Table[cell[0]][cell[1]] {
			if entity == e {
				inCell = true
			}
		}
		if !inCell {
			t.Fatal(fmt.Sprintf("%v,%v was not mapped to cell %v",
				e.GetVec2D("Position"),
				e.GetVec2D("Box"),
				cell))
		}
	}
}

func TestSpatialHashTableCopy(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	testingSpawnSpatial(w, Vec2D{1, 1}, Vec2D{1, 1})
	w.Update(FRAME_SLEEP_MS / 2)
	table := sh.Table
	tableCopy := sh.TableCopy()
	if table[0][0][0] != tableCopy[0][0][0] {
		t.Fatal("CurrentTableCopy() doesn't return a copy")
	}
	table[0][0] = table[0][0][:0]
	if len(tableCopy[0][0]) == 0 {
		t.Fatal("CurrentTableCopy() doesn't return a copy")
	}
}

func TestSpatialHashTableToString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	w := testingWorld()
	sh := NewSpatialHashSystem(10, 10)
	w.AddSystems(sh)
	table := sh.Table
	s0 := table.String()
	for i := 0; i < 500; i++ {
		testingSpawnSpatial(w,
			Vec2D{rand.Float64() * 1024, rand.Float64() * 1024},
			Vec2D{5, 5})
	}
	w.Update(FRAME_SLEEP_MS)
	s1 := sh.Table.String()
	if len(s1) < len(s0) {
		t.Fatal("spatial hash did not show entities in its String() representation")
	}
}
