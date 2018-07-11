package engine

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"go.uber.org/atomic"
)

// used to store an entity with a position in a grid cell
type EntityGridPosition struct {
	entity *EntityToken
	x, y   int
}

// the actual cell data structure is a gridDimension x gridDimension array of []EntityGridPosition
type SpatialHashTable [][][]EntityGridPosition

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
	// channels (one per grid square) used to receive successfully-read
	// positions from the goroutines which scan the entities
	// (not double-buffered since we exclude two
	// Update() instances from running at the same time using
	// the computeInProgress flag)
	cellChannels [][]chan EntityGridPosition
	// used to signal that the compute is done from one of the cellReceiver()
	// workers
	computeDoneChannel chan bool
	// used for the cellReceiver() worker to signal it is ending, because it
	// noticed there are no more entities to scan in the current computation
	cellReceiverStopChannel chan bool
	// used for the computation to wait for all cell receivers to end before
	// allowing a new computation to start
	cellReceiverWaitGroup sync.WaitGroup
	// how many entities are yet to store in the current computation
	entitiesRemaining atomic.Uint32
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
		cellReceiverStopChannel: make(chan bool),
		computeDoneChannel:      make(chan bool),
	}
	h.gridX = gridX
	h.gridY = gridY
	h.tables[0] = make([][][]EntityGridPosition, gridX)
	h.tables[1] = make([][][]EntityGridPosition, gridX)
	h.cellChannels = make([][]chan EntityGridPosition, gridX)
	// for each column (x) in the grid
	for x := 0; x < gridX; x++ {
		h.tables[0][x] = make([][]EntityGridPosition, gridY)
		h.tables[1][x] = make([][]EntityGridPosition, gridY)
		h.cellChannels[x] = make([]chan EntityGridPosition, gridY)
		// for each cell in the row (y)
		for y := 0; y < gridY; y++ {
			h.tables[0][x][y] = make([]EntityGridPosition, MAX_ENTITIES)
			h.tables[1][x][y] = make([]EntityGridPosition, MAX_ENTITIES)
			h.cellChannels[x][y] = make(chan EntityGridPosition, MAX_ENTITIES)

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

	// set the computingTable pointer used by the cellReceiver() workers to point
	// to the table we're building
	h.computingTable = &h.tables[h.nextBufIndex()]

	// set the number of entities remaining (used by cellReceiver workers to
	// notify that computation is done)
	h.entitiesRemaining.Store(uint32(len(h.spatialEntities.Entities)))
	// clear any old data and run the computation
	h.clearComputingTable()
	h.startCellReceivers()
	h.startEntityScanners()
	// wait for computation to finish
	<-h.computeDoneChannel
	// if we're here, the computation has completed.
	// this increment, due to the modulo logic, is equivalent to setting
	// computedBufIndex = nextBufIndex
	h.tableBufIndex.Inc()
	h.cellReceiverWaitGroup.Wait()
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

func (h *SpatialHashSystem) startCellReceivers() {
	// start the cellReceiver for this cell
	for x := 0; x < h.gridX; x++ {
		for y := 0; y < h.gridY; y++ {
			h.cellReceiverWaitGroup.Add(1)
			go h.cellReceiver(x, y, h.cellChannels[x][y])
		}
	}
}

// used to receive EntityGridPositions and put them into the right cell
func (h *SpatialHashSystem) cellReceiver(x int, y int, c chan EntityGridPosition) {
	defer h.cellReceiverWaitGroup.Done()
	for {
		select {
		case _ = <-h.cellReceiverStopChannel:
			return
		case entityPosition := <-c:
			h.storeEntity(x, y, entityPosition)
			if h.entitiesRemaining.Dec() == 0 {
				for i := 0; i < (h.gridX*h.gridY)-1; i++ {
					h.cellReceiverStopChannel <- true
				}
				h.computeDoneChannel <- true
				return
			}
		}
	}
}

func (h *SpatialHashSystem) storeEntity(
	x int, y int, entityPosition EntityGridPosition) {

	cell := &((*h.computingTable)[x][y])
	*cell = append(*cell, entityPosition)
}

func (h *SpatialHashSystem) startEntityScanners() {
	// Divide the list of entities into a certain number of partitions
	// which will be scanned by entityScanner() instances.
	nScanPartitions := runtime.NumCPU()
	partition_size := len(h.spatialEntities.Entities) / nScanPartitions
	// only compute if there is at least 1 entity
	if len(h.spatialEntities.Entities) > 0 {
		for partition := 0; partition < nScanPartitions; partition++ {
			offset := partition * partition_size
			if partition == nScanPartitions-1 {
				// the last partition includes the remainder
				partition_size = len(h.spatialEntities.Entities) - offset
			}
			go h.entityScanner(offset, partition_size)
		}
	}
}

// used to iterate the entities and send them to the right cell's channels
func (h *SpatialHashSystem) entityScanner(offset int, partition_size int) {
	// iterate each entity in the partition
	for i := 0; i < partition_size; i++ {
		// get the entity's box
		entity := h.spatialEntities.Entities[offset+i]
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
				h.cellChannels[x][y] <- EntityGridPosition{entity, x, y}
			}
		}
	}
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
	t2 := make([][][]EntityGridPosition, h.gridX)
	for x := 0; x < h.gridX; x++ {
		t2[x] = make([][]EntityGridPosition, h.gridX)
		for y := 0; y < h.gridY; y++ {
			t2[x][y] = make([]EntityGridPosition, len((*t)[x][y]))
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
			size += int(unsafe.Sizeof(EntityGridPosition{})) * len(cell)
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
