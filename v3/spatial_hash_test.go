package sameriver

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestSpatialHashInsertion(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
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
		e := testingSpawnSpatial(w, posbox[0], posbox[1])
		entityCells[e] = cells
	}
	w.Update(FRAME_DURATION_INT / 2)
	for e, cells := range entityCells {
		for _, cell := range cells {
			inCell := false
			for _, entity := range sh.Hasher.Entities(cell[0], cell[1]) {
				if entity == e {
					inCell = true
				}
			}
			if !inCell {
				t.Fatalf("%v,%v was not mapped to cell %v",
					e.GetVec2D("Position"),
					e.GetVec2D("Box"),
					cell)
			}
		}
	}
}

func TestSpatialHashMany(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	for i := 0; i < 300; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	w.Update(FRAME_DURATION_INT / 2)
	n_entities := len(w.GetActiveEntitiesSet())
	seen := make(map[*Entity]bool)
	found := 0
	table := sh.Hasher.TableCopy()
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			cell := table[x][y]
			for _, e := range cell {
				if _, ok := seen[e]; !ok {
					found++
					seen[e] = true
				}
			}
		}
	}
	if found != n_entities {
		t.Fatal("Some entities were not in any cell!")
	}
}

func TestSpatialHashLargeEntity(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	pos := Vec2D{20, 20}
	box := Vec2D{5, 5}
	cells := [][2]int{
		[2]int{1, 1},
		[2]int{1, 2},
		[2]int{2, 1},
		[2]int{2, 2},
	}
	e := testingSpawnSpatial(w, pos, box)
	w.Update(FRAME_DURATION_INT / 2)
	for _, cell := range cells {
		inCell := false
		for _, entity := range sh.Hasher.Entities(cell[0], cell[1]) {
			if entity == e {
				inCell = true
			}
		}
		if !inCell {
			t.Fatalf("%v,%v was not mapped to cell %v",
				e.GetVec2D("Position"),
				e.GetVec2D("Box"),
				cell)
		}
	}
}

func TestSpatialHashCellsWithinDistance(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)

	box := Vec2D{0, 0}

	// we're checking the radius at 0, 0, the corner of the world
	cells := sh.Hasher.CellsWithinDistance(Vec2D{0, 0}, box, 25.0)
	if len(cells) != 8 {
		t.Fatalf("circle centered at 0, 0 of radius 25 should touch 8 cells; got %d: %v", len(cells), cells)
	}
	cells = sh.Hasher.CellsWithinDistance(Vec2D{0, 0}, box, 29.0)
	if len(cells) != 9 {
		t.Fatalf("circle centered at 0, 0 of radius 29 should touch 9 cells; got %d: %v", len(cells), cells)
	}
	// now check from a position not quite at the corner
	cells = sh.Hasher.CellsWithinDistance(Vec2D{20, 20}, box, 29.0)
	if len(cells) != 25 {
		t.Fatalf("circle centered at 20, 20 of radius 29 should touch 25 cells; got %d: %v", len(cells), cells)
	}
	cells = sh.Hasher.CellsWithinDistance(Vec2D{20, 20}, box, 7.0)
	if len(cells) != 4 {
		t.Fatalf("circle centered at 20, 20 of radius 7 should touch 4 cells; got %d: %v", len(cells), cells)
	}
	cells = sh.Hasher.CellsWithinDistance(Vec2D{25, 25}, box, 1.0)
	if len(cells) != 1 {
		t.Fatalf("circle centered at 25, 25 of radius 1 should touch 1 cell; got %d: %v", len(cells), cells)
	}
}

func TestSpatialHashEntitiesWithinDistance(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)

	e := testingSpawnSpatial(w, Vec2D{50, 50}, Vec2D{5, 5})

	near := make([]*Entity, 0)
	far := make([]*Entity, 0)
	for spawnRadius := 30.0; spawnRadius <= 38; spawnRadius += 8 {
		for i := 0.0; i < 360; i += 10 {
			theta := 2.0 * math.Pi * (i / 360)
			offset := Vec2D{
				spawnRadius * math.Cos(theta),
				spawnRadius * math.Sin(theta),
			}
			spawned := testingSpawnSpatial(w,
				e.GetVec2D("Position").Add(offset),
				Vec2D{5, 5})
			if spawnRadius == 30.0 {
				near = append(near, spawned)
			} else {
				far = append(far, spawned)
			}
		}
	}
	w.Update(FRAME_DURATION_INT / 2)
	nearGot := w.EntitiesWithinDistance(
		*e.GetVec2D("Position"),
		*e.GetVec2D("Box"),
		30.0)
	if len(nearGot) != 37 {
		t.Fatalf("Should be 37 near entities; got %d", len(nearGot))
	}
	for _, eNear := range near {
		if indexOfEntityInSlice(&nearGot, eNear) == -1 {
			t.Fatal("Did not find expected near entity")
		}
	}
	for _, eFar := range far {
		if indexOfEntityInSlice(&nearGot, eFar) != -1 {
			t.Fatal("Shouldn't have returned far entity")
		}
	}
}

func TestSpatialHashEntitiesWithinDistanceApprox(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})

	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)

	e := testingSpawnSpatial(w, Vec2D{50, 50}, Vec2D{5, 5})

	near := make([]*Entity, 0)
	for spawnRadius := 30.0; spawnRadius <= 38; spawnRadius += 8 {
		for i := 0.0; i < 360; i += 10 {
			theta := 2.0 * math.Pi * (i / 360)
			offset := Vec2D{
				spawnRadius * math.Cos(theta),
				spawnRadius * math.Sin(theta),
			}
			spawned := testingSpawnSpatial(w,
				e.GetVec2D("Position").Add(offset),
				Vec2D{5, 5})
			if spawnRadius == 30.0 {
				near = append(near, spawned)
			}
		}
	}
	w.Update(FRAME_DURATION_INT / 2)
	nearGot := sh.Hasher.EntitiesWithinDistanceApprox(
		*e.GetVec2D("Position"),
		*e.GetVec2D("Box"),
		30.0)
	// somehow by coincidence - and this can only be a sign that god exists -
	// the number of near entities by exact calculation is 37, but the number
	// of near entities by approx calculation (with the grid we defined) is 73.
	// absolutely beautiful. ... well when you think about it, it makes sense.
	// 36 + 1 = 37, 2*36 + 1 = 73
	// but, that's still beautiful, and still a sign.
	if len(nearGot) != 73 {
		t.Fatalf("Should be 73 near entities; got %d", len(nearGot))
	}
	for _, eNear := range near {
		if indexOfEntityInSlice(&nearGot, eNear) == -1 {
			t.Fatal("Did not find expected near entity")
		}
	}
}

func TestSpatialHashTableCopy(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	testingSpawnSpatial(w, Vec2D{1, 1}, Vec2D{1, 1})
	w.Update(FRAME_DURATION_INT / 2)
	w.Update(FRAME_DURATION_INT / 2)
	table := sh.Hasher.Table
	tableCopy := sh.Hasher.TableCopy()
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
	w.RegisterSystems(sh)
	s0 := sh.Hasher.String()
	for i := 0; i < 500; i++ {
		testingSpawnSpatial(w,
			Vec2D{rand.Float64() * 1024, rand.Float64() * 1024},
			Vec2D{5, 5})
	}
	w.Update(FRAME_DURATION_INT)
	s1 := sh.Hasher.String()
	if len(s1) < len(s0) {
		t.Fatal("spatial hash did not show entities in its String() representation")
	}
}

func TestSpatialHashExpand(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	testData := map[[2]Vec2D][][2]int{
		[2]Vec2D{Vec2D{5, 5}, Vec2D{10, 10}}: [][2]int{
			[2]int{0, 0},
			[2]int{0, 1},
			[2]int{1, 0},
			[2]int{1, 1},
		},
	}
	for posbox := range testData {
		testingSpawnSpatial(w, posbox[0], posbox[1])
	}
	oldCapacity := cap(sh.Hasher.Table[0][0])
	Logger.Printf("oldCapacity: %d", oldCapacity)

	w.Update(FRAME_DURATION_INT / 2)
	Logger.Println(sh.Hasher.Table)

	sh.Expand(MAX_ENTITIES / 2)

	if !(oldCapacity < cap(sh.Hasher.Table[0][0])) {
		t.Fatal("Did not expand capacity of cells")
	}

	expected := [][3]int{
		[3]int{0, 0, 1},
		[3]int{0, 1, 1},
		[3]int{1, 0, 1},
		[3]int{1, 1, 1},
	}
	for _, e := range expected {
		x, y, n := e[0], e[1], e[2]
		if len(sh.Hasher.Table[x][y]) != n {
			Logger.Printf("[%d][%d]", x, y)
			Logger.Println(sh.Hasher.Table[x][y])
			t.Fatal("altered cell counts")
		}
	}
}
