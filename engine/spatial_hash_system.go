package engine

import (
	"bytes"
	"fmt"
	"unsafe"

	"go.uber.org/atomic"
)

// the actual cell data structure is a gridDimension x gridDimension array of
// entities
type SpatialHashTable [][][]*EntityToken

// used to compute the spatial hash tables given a list of entities
type SpatialHashSystem struct {
	// the world entities live in
	w *World
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	spatialEntities *UpdatedEntityList
	// basic data members needed to divide the world into cells and
	// store the entity data in each cell
	gridX int
	gridY int
	// we double-buffer the data structure by atomically incrementing this
	// Uint32 every time we finish computing a new one, taking its
	// (value - 1) % 2 to return the index of the latest completed cell
	// structure
	tableBufIndex atomic.Uint32
	// tables is a double-buffer for holding SpatialHashTable results
	// computed during Update
	tables [2]SpatialHashTable
	// an unsafe pointer used by cellReceiver() workers in the compute stage
	// to find the table we're currently building
	computingTable *SpatialHashTable
	computedTable  *SpatialHashTable
	// a lock to ensure that we never enter Update while
	// another instance of the function is still running (we return early,
	// while a mutex would freeze the goroutine of the caller if called
	// in sync, or at least lead to leaked goroutines stacking up if
	// each call was spawned in a goroutine and was consistently failing
	// to execute before the next call)
	ComputeRunning atomic.Uint32
}

func NewSpatialHashSystem(gridX int, gridY int) *SpatialHashSystem {
	h := SpatialHashSystem{
		gridX: gridX,
		gridY: gridY,
	}
	h.tables[0] = make([][][]*EntityToken, gridX)
	h.tables[1] = make([][][]*EntityToken, gridX)
	// for each column (x) in the grid
	for x := 0; x < gridX; x++ {
		h.tables[0][x] = make([][]*EntityToken, gridY)
		h.tables[1][x] = make([][]*EntityToken, gridY)
		// for each cell in the row (y)
		for y := 0; y < gridY; y++ {
			h.tables[0][x][y] = make([]*EntityToken, MAX_ENTITIES)
			h.tables[1][x][y] = make([]*EntityToken, MAX_ENTITIES)
		}
	}
	return &h
}

func (s *SpatialHashSystem) LinkWorld(w *World) {
	s.w = w
	// get a list of spatial entities
	s.spatialEntities = w.em.GetUpdatedEntityList(
		EntityQueryFromComponentBitArray("spatial",
			MakeComponentBitArray(
				[]ComponentType{POSITION_COMPONENT, BOX_COMPONENT})))
}

// spawns a certain number of goroutines to iterate through entities, trying
// to lock them and get their position and send the entities and their
// positions to another set of goroutines handling the building of the
// list for each grid cell
func (h *SpatialHashSystem) Update(dt_ms float64) {

	// this lock prevents another call to Update()
	// entering while we are currently calculating (this ensures robustness
	// if for some reason it is called too often)
	if !h.ComputeRunning.CAS(0, 1) {
		return
	}

	entities := h.spatialEntities.entities
	// we will write to the table indicated by nextBufIndex
	h.computingTable = &h.tables[h.nextBufIndex()]
	// clear any old data and run the computation
	h.clearComputingTable()
	h.scanAndInsertEntities(entities)
	// this increment, due to the modulo 2 logic, is equivalent to setting
	// computedBufIndex = nextBufIndex
	h.tableBufIndex.Inc()
	h.ComputeRunning.Store(0)
}

func (h *SpatialHashSystem) clearComputingTable() {
	// NOTE: we "clear" the slice by setting its length to 0 (capacity remains
	// allocated, this will cause a negligible memory "waste" if entities
	// cluster in a cell but never somewhere else. Maybe this could matter if
	// MAX_ENTITIES eventually clustered in each cell, but that's unlikely)
	for x := 0; x < h.gridX; x++ {
		for y := 0; y < h.gridY; y++ {
			cell := &((*h.computingTable)[x][y])
			*cell = (*cell)[:0]
		}
	}
}

// used to iterate the entities and send them to the right cells
func (h *SpatialHashSystem) scanAndInsertEntities(entities []*EntityToken) {

	// iterate each entity in the partition
	for _, entity := range entities {
		pos := h.w.em.Components.Position[entity.ID]
		box := h.w.em.Components.Box[entity.ID]
		// find out how many grids the entity spans in x and y (almost always 0,
		// but we want to be thorough, and the fact that it's got a predictable
		// pattern 99% of the time means that branch prediction should help us)
		gridsWide := int(box.X) / (h.w.Height / h.gridX)
		gridsHigh := int(box.Y) / (h.w.Width / h.gridY)
		// figure out which cell the topleft corner is in
		topLeftCellX, topLeftCellY := h.cellForPoint(int(pos.X), int(pos.Y))
		// walk through each cell the entity touches by starting in the top-left
		// and walking according to gridsHigh and gridsWide
		for ix := 0; ix < gridsWide+1; ix++ {
			for iy := 0; iy < gridsHigh+1; iy++ {
				x := topLeftCellX + ix
				y := topLeftCellY + iy
				h.storeEntity(x, y, entity)
			}
		}
	}
}

func (h *SpatialHashSystem) storeEntity(x int, y int, entity *EntityToken) {
	cell := &((*h.computingTable)[x][y])
	*cell = append(*cell, entity)
}

// helper function for hashing
func (h *SpatialHashSystem) cellForPoint(x int, y int) (cellX int, cellY int) {
	return x / (h.w.Width / h.gridX), y / (h.w.Height / h.gridY)
}

// returns the double-buffer index of the current spatial hash table (the one
// that callers will want). We subtract 1 since we increment tableBufIndex
// atomically once the computation completes
func (h *SpatialHashSystem) computedBufIndex() uint32 {
	return (h.tableBufIndex.Load() - 1) % 2
}

// returns the double-buffer index of the *next* spatial hash table (the one
// we will be / we are computing)
func (h *SpatialHashSystem) nextBufIndex() uint32 {
	return h.tableBufIndex.Load() % 2
}

// get the pointer to the current spatial hash data structure
// NOTE: this pointer is only safe to use until you've called Update
// at *most* one time more. If you can't ensure that it won't be called, and
// want to do something outside of the main game loop with a spatial hash
// result, use CurrentTableCopy()
func (h *SpatialHashSystem) CurrentTablePointer() *SpatialHashTable {
	return &h.tables[h.computedBufIndex()]
}

// get a *copy* of the current table which is safe to hold onto, mutate, etc.
func (h *SpatialHashSystem) CurrentTableCopy() SpatialHashTable {
	var t = h.CurrentTablePointer()
	t2 := make([][][]*EntityToken, h.gridX)
	for x := 0; x < h.gridX; x++ {
		t2[x] = make([][]*EntityToken, h.gridX)
		for y := 0; y < h.gridY; y++ {
			t2[x][y] = make([]*EntityToken, len((*t)[x][y]))
			copy(t2[x][y], (*t)[x][y])

		}
	}
	return t2
}

// turn a SpatialHashTable into a String representation (NOTE: do *NOT* call
// this on a pointer returned from CurrentTablePointer unless you can be sure
// that you have not called Update more than once - it does not
// lock the table it reads, and if you call Update twice, you may
// start to write to the table as this function reads it). Usually best to call
// on a Copy()
func (h *SpatialHashSystem) String(table *SpatialHashTable) string {
	var buffer bytes.Buffer
	size := int(unsafe.Sizeof(*table))
	buffer.WriteString("[\n")
	for x := 0; x < h.gridX; x++ {
		for y := 0; y < h.gridY; y++ {
			cell := (*table)[x][y]
			size += int(unsafe.Sizeof(&EntityToken{})) * len(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				fmt.Sprintf("%+v", cell)))
			if !(y == h.gridY-1 && x == h.gridX-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}
