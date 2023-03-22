package sameriver

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"sync"
	"unsafe"
)

type SpatialHasher struct {
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	SpatialEntities *UpdatedEntityList
	// basic data members needed to divide the world into cells and
	// store the entity data in each cell
	GridX     int
	GridY     int
	CellSizeX float64
	CellSizeY float64
	// table of cells, GridX x GridY, that holds the entities
	Table [][][]*Entity

	// used in scanAndInsertEntitiesparallelC
	tableMutexes [][]sync.Mutex

	// capacity keeps track of the world's max entities
	// so we can keep the right capacity (max entities / 4) in each grid cell
	capacity int
}

func NewSpatialHasher(gridX, gridY int, w *World) *SpatialHasher {
	h := &SpatialHasher{
		GridX:     gridX,
		GridY:     gridY,
		CellSizeX: w.Width / float64(gridX),
		CellSizeY: w.Height / float64(gridY),
		capacity:  w.MaxEntities(),
	}
	h.allocTable()
	h.allocTableMutexes()
	// get spatial entities from world
	h.SpatialEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray("spatial",
			w.em.components.BitArrayFromIDs([]ComponentID{POSITION, BOX})))

	return h
}

func (h *SpatialHasher) allocTable() {
	h.Table = make([][][]*Entity, h.GridX)
	// for each column (x)
	for x := 0; x < h.GridX; x++ {
		h.Table[x] = make([][]*Entity, h.GridY)
		// for each cell in the row (y)
		for y := 0; y < h.GridY; y++ {
			h.Table[x][y] = make([]*Entity, 0, h.capacity/4)
		}
	}
}

func (h *SpatialHasher) allocTableMutexes() {
	h.tableMutexes = make([][]sync.Mutex, h.GridY)
	for x := 0; x < h.GridX; x++ {
		h.tableMutexes[x] = make([]sync.Mutex, h.GridX)
	}
}

func (h *SpatialHasher) Entities(x, y int) []*Entity {
	return h.Table[x][y]
}

func (h *SpatialHasher) Update() {
	// if we only have 1 CPU, use single-threaded (don't needlessly use mutexes)
	// otherwise, single vs parallel isn't exactly clear which is better
	// (see benchmark_spatial_hash_compare.sh); it depends on grid size and current CPU
	// load. Let's assume all things being equal that parallel will be better if we
	// have the cores for it
	if runtime.NumCPU() == 1 {
		h.singleThreadUpdate()
	} else {
		h.parallelUpdateC()
	}
}

func (h *SpatialHasher) parallelUpdateC() {
	h.clearTable()
	h.scanAndInsertEntitiesparallelC()
}
func (h *SpatialHasher) singleThreadUpdate() {
	h.clearTable()
	h.scanAndInsertEntitiesSingleThread()
}

func (h *SpatialHasher) clearTable() {
	// NOTE: we "clear" the slice by setting its length to 0 (capacity remains
	// allocated, this will cause a negligible memory "waste" if entities
	// cluster in a cell but never somewhere else. Maybe this could matter if
	// MAX_ENTITIES eventually clustered in each cell, but that's unlikely)
	for x := 0; x < h.GridX; x++ {
		for y := 0; y < h.GridY; y++ {
			cell := &h.Table[x][y]
			*cell = (*cell)[:0]
		}
	}
}

func (h *SpatialHasher) CellRangeOfRect(pos, box Vec2D) (cellX0, cellX1, cellY0, cellY1 int) {
	cellX0 = int(pos.X / h.CellSizeX)
	cellX1 = int((pos.X + box.X) / h.CellSizeX)
	cellY0 = int(pos.Y / h.CellSizeY)
	cellY1 = int((pos.Y + box.Y) / h.CellSizeY)
	return cellX0, cellX1, cellY0, cellY1
}

// 72912 ns/op (at GridX,GridY = 10,10)
// ^^^ this performance was only seen under certain conditions. More often
// we get something much closer to single-threaded performance
func (h *SpatialHasher) scanAndInsertEntitiesparallelC() {
	numWorkers := runtime.NumCPU()
	// Launch workers to scan and insert into their own tables
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			startIdx := int(math.Floor(float64(len(h.SpatialEntities.entities)) * float64(workerID) / float64(numWorkers)))
			endIdx := int(math.Floor(float64(len(h.SpatialEntities.entities)) * float64(workerID+1) / float64(numWorkers)))

			for j := startIdx; j < endIdx; j++ {
				e := h.SpatialEntities.entities[j]
				pos := e.GetVec2D(POSITION)
				box := e.GetVec2D(BOX)
				cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(pos.ShiftedCenterToBottomLeft(*box), *box)

				for y := cellY0; y <= cellY1; y++ {
					for x := cellX0; x <= cellX1; x++ {
						if x < 0 || x >= h.GridX || y < 0 || y >= h.GridY {
							continue
						}
						h.tableMutexes[x][y].Lock()
						h.Table[x][y] = append(h.Table[x][y], e)
						h.tableMutexes[x][y].Unlock()
					}
				}
			}
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()
}

// 104519 ns/op (at GridX,GridY = 10,10)
// somewhat suprisingly, better than some parallel versions
func (h *SpatialHasher) scanAndInsertEntitiesSingleThread() {
	for _, e := range h.SpatialEntities.entities {
		pos := e.GetVec2D(POSITION)
		box := e.GetVec2D(BOX)

		// walk through each cell the entity touches by
		// starting in the bottom-left and walking cell by cell
		// through each row to the top-right
		cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(pos.ShiftedCenterToBottomLeft(*box), *box)
		for x := cellX0; x <= cellX1; x++ {
			for y := cellY0; y <= cellY1; y++ {
				if x < 0 || x >= h.GridX || y < 0 || y >= h.GridY {
					continue
				}
				cell := &h.Table[x][y]
				*cell = append(*cell, e)
			}
		}
	}
}

// TableCopy gets a *copy* of the current table which is safe to hold onto, mutate, etc.
func (h *SpatialHasher) TableCopy() [][][]*Entity {
	t2 := make([][][]*Entity, h.GridX)
	for x := 0; x < h.GridX; x++ {
		t2[x] = make([][]*Entity, h.GridX)
		for y := 0; y < h.GridY; y++ {
			t2[x][y] = make([]*Entity, len(h.Table[x][y]))
			copy(t2[x][y], h.Table[x][y])
		}
	}
	return t2
}

func (h *SpatialHasher) GetCellPosAndBox(x, y int) (pos, box Vec2D) {
	pos = Vec2D{
		float64(x) * h.CellSizeX,
		float64(y) * h.CellSizeY}
	box = Vec2D{
		h.CellSizeX,
		h.CellSizeY}
	return pos, box
}

// CellsWithinDistance calculates which cells are within the distance d from the closest
// point on the box centered at pos (imagine a rounded-corner box
// extending d past the limits of the box)
func (h *SpatialHasher) CellsWithinDistance(pos, box Vec2D, d float64) [][2]int {
	cells := make([][2]int, 0)

	// first approximate which cells might be valid by simply
	// extending the box by +d in each direction
	approximatorPos := pos.ShiftedCenterToBottomLeft(box).Sub(Vec2D{d, d})
	approximatorBox := Vec2D{2*d + box.X, 2*d + box.Y}
	cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(approximatorPos, approximatorBox)
	candidateCells := make([][2]int, 0)
	for x := cellX0; x <= cellX1; x++ {
		for y := cellY0; y <= cellY1; y++ {
			if x < 0 || x >= h.GridX || y < 0 || y >= h.GridY {
				continue
			}
			candidateCells = append(candidateCells, [2]int{x, y})
		}
	}
	for _, cellXY := range candidateCells {
		cellPos, cellBox := h.GetCellPosAndBox(cellXY[0], cellXY[1])
		if RectWithinDistanceOfRect(cellPos, cellBox, pos.ShiftedCenterToBottomLeft(box), box, d) {
			cells = append(cells, cellXY)
		}
	}
	return cells
}

// CellsWithinDistanceApprox extends the box +d on all sides and returns the cells it touches
// (NOTE: the corners will slightly over-estimate since they should
// truly be rounded)
// but it's a faster calculation
func (h *SpatialHasher) CellsWithinDistanceApprox(pos, box Vec2D, d float64) [][2]int {
	approximatorPos := pos.ShiftedCenterToBottomLeft(box).Sub(Vec2D{d, d})
	approximatorBox := Vec2D{2*d + box.X, 2*d + box.Y}
	cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(approximatorPos, approximatorBox)
	candidateCells := make([][2]int, 0)
	for x := cellX0; x <= cellX1; x++ {
		for y := cellY0; y <= cellY1; y++ {
			if x < 0 || x >= h.GridX || y < 0 || y >= h.GridY {
				continue
			}
			candidateCells = append(candidateCells, [2]int{x, y})
		}
	}
	return candidateCells
}

func (h *SpatialHasher) EntitiesWithinDistanceApprox(
	pos, box Vec2D, d float64) []*Entity {
	return h.EntitiesWithinDistanceApproxFilter(pos, box, d,
		func(e *Entity) bool { return true })
}

// EntitiesWithinDistanceApproxFilter uses the approx distance since it's faster.
// Overestimates slightly diagonally.
func (h *SpatialHasher) EntitiesWithinDistanceApproxFilter(
	pos, box Vec2D, d float64, predicate func(*Entity) bool) []*Entity {
	results := make([]*Entity, 0)
	found := make(map[int]*Entity)
	cells := h.CellsWithinDistanceApprox(pos, box, d)
	// for each cell, note that we found the entity in the map
	// initial naive slice append was quadruple-counting entities that
	// sat on the intersection of four cells, etc.
	for _, cell := range cells {
		x := cell[0]
		y := cell[1]
		for _, e := range h.Table[x][y] {
			found[e.ID] = e
		}
	}
	for _, e := range found {
		if predicate(e) {
			results = append(results, e)
		}
	}
	return results
}

// NOTE: can return inactive entities
func (h *SpatialHasher) EntitiesWithinDistance(pos, box Vec2D, d float64) []*Entity {
	return h.EntitiesWithinDistanceFilter(pos, box, d,
		func(e *Entity) bool { return true })
}

// NOTE: can return inactive entitites
func (h *SpatialHasher) EntitiesWithinDistanceFilter(
	pos, box Vec2D, d float64, predicate func(*Entity) bool) []*Entity {
	candidates := h.EntitiesWithinDistanceApprox(pos, box, d)
	results := make([]*Entity, 0)
	for _, e := range candidates {
		ePos := *e.GetVec2D(POSITION)
		eBox := *e.GetVec2D(BOX)
		if predicate(e) && RectWithinDistanceOfRect(
			pos.ShiftedCenterToBottomLeft(box), box,
			ePos.ShiftedCenterToBottomLeft(eBox), eBox,
			d) {
			results = append(results, e)
		}
	}
	return results
}

// String turns a SpatialHashTable into a String representation (NOTE: do *NOT* call
// this on a pointer returned from CurrentTablePointer unless you can be sure
// that you have not called Update more than once - it does not
// lock the table it reads, and if you call Update twice, you may
// start to write to the table as this function reads it). Usually best to call
// on a Copy()
func (h *SpatialHasher) String() string {
	var buffer bytes.Buffer
	size := int(unsafe.Sizeof(h.Table))
	buffer.WriteString("[\n")
	for x := 0; x < h.GridX; x++ {
		for y := 0; y < h.GridY; y++ {
			cell := h.Table[x][y]
			size += int(unsafe.Sizeof(&Entity{})) * cap(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				EntitySliceToString(cell)))
			if !(y == h.GridY-1 && x == h.GridX-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}

func (h *SpatialHasher) Expand(n int) {
	for x := 0; x < h.GridX; x++ {
		for y := 0; y < h.GridY; y++ {
			oldCount := len(h.Table[x][y])
			newCell := make([]*Entity, (h.capacity+n)/4)
			copy(newCell, h.Table[x][y])
			h.Table[x][y] = newCell[:oldCount]
		}
	}
	h.capacity += n
}
