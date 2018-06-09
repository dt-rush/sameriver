package engine

import (
	"bytes"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"go.uber.org/atomic"
)

// used to store an entity with a position in a grid cell
type EntityPosition struct {
	entity   EntityToken
	position sdl.Rect
}

// the actual cell data structure is a GRID x GRID array of []EntityPosition
type SpatialHashTable [][][]EntityPosition

// used to compute the spatial hash tables given a list of entities
type SpatialHash struct {
	// basic data members needed to divide the world into cells and
	// store the entity data in each cell
	WORLD_WIDTH  int
	WORLD_HEIGHT int
	GRID         int
	// we double-buffer the data structure by atomically incrementing this
	// Uint32 every time we finish computing a new one, taking its
	// (value - 1) % 2 to return the index of the latest completed cell
	// structure
	tableBufIndex atomic.Uint32
	// tables is a double-buffer for holding SpatialHashTable results
	// computed during ComputeSpatialHash
	tables [2]SpatialHashTable
	// an unsafe pointer used by receiver() workers in the compute stage
	// to find the table we're currently building
	computingTable *SpatialHashTable
	computedTable  *SpatialHashTable
	// channels (one per grid square) used to receive successfully-read
	// positions from the goroutines which scan the entities
	// (not double-buffered since we exclude two
	// ComputeSpatialHash() instances from running at the same time using
	// the computeInProgress flag)
	cellChannels [][]chan EntityPosition
	// used to signal that the compute is done from one of the receive()
	// workers
	computeDoneChannel chan bool
	// how many entities are yet to store in the current computation
	entitiesRemaining atomic.Uint32
	// a lock to ensure that we never enter ComputeSpatialHash while
	// another instance of the function is still running (we return early,
	// while a mutex would freeze the goroutine of the caller if called
	// in sync, or at least lead to leaked goroutines stacking up if
	// each call was spawned in a goroutine and was consistently failing
	// to execute before the next call)
	canEnterCompute atomic.Uint32
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	spatialEntities *UpdatedEntityList
	// entityManager is used to acquire locks on entities
	em *EntityManager
}

func NewSpatialHash(
	WORLD_WIDTH int,
	WORLD_HEIGHT int,
	GRID int,
	em *EntityManager) *SpatialHash {

	h := SpatialHash{computeDoneChannel: make(chan bool)}
	h.WORLD_WIDTH = WORLD_WIDTH
	h.WORLD_HEIGHT = WORLD_HEIGHT
	h.GRID = GRID
	h.tables[0] = make([][][]EntityPosition, GRID)
	h.tables[1] = make([][][]EntityPosition, GRID)
	h.cellChannels = make([][]chan EntityPosition, GRID)
	// for each row in the grid
	for y := 0; y < GRID; y++ {
		h.tables[0][y] = make([][]EntityPosition, GRID)
		h.tables[1][y] = make([][]EntityPosition, GRID)
		h.cellChannels[y] = make([]chan EntityPosition, GRID)
		// for each cell in the row
		for x := 0; x < GRID; x++ {
			h.tables[0][y][x] = make([]EntityPosition, MAX_ENTITIES)
			h.tables[1][y][x] = make([]EntityPosition, MAX_ENTITIES)
			h.cellChannels[y][x] = make(chan EntityPosition, MAX_ENTITIES)
			// start the receiver for this cell
			go h.receiver(y, x, h.cellChannels[y][x])
		}
	}
	// get a list of spatial entities
	h.spatialEntities = em.GetUpdatedEntityList(
		EntityQueryFromComponentBitArray("spatial",
			MakeComponentBitArray([]ComponentType{BOX_COMPONENT})))
	// take down references needed during compute
	h.em = em
	// canEnterCompute is expressed in the traditional sense that 1 = true,
	// so we have to initialize it
	h.canEnterCompute.Store(1)
	return &h
}

// returns the double-buffer index of the current spatial hash table (the one
// that callers will want). We subtract 1 since we increment tableBufIndex
// atomically once the computation completes
func (h *SpatialHash) computedBufIndex() uint32 {
	return (h.tableBufIndex.Load() - 1) % 2
}

// returns the double-buffer index of the *next* spatial hash table (the one
// we will be / we are computing)
func (h *SpatialHash) nextBufIndex() uint32 {
	return h.tableBufIndex.Load() % 2
}

// get the pointer to the current spatial hash data structure
// NOTE: this pointer is only safe to use until you've called ComputeSpatialHash
// at *most* one time more. If you can't ensure that it won't be called, and
// want to do something outside of the main game loop with a spatial hash
// result, use CurrentTableCopy()
func (h *SpatialHash) CurrentTablePointer() *SpatialHashTable {
	return &h.tables[h.computedBufIndex()]
}

// get a *copy* of the current table which is safe to hold onto, mutate, etc.
func (h *SpatialHash) CurrentTableCopy() SpatialHashTable {
	var t = h.CurrentTablePointer()
	t2 := make([][][]EntityPosition, h.GRID)
	for y := 0; y < h.GRID; y++ {
		t2[y] = make([][]EntityPosition, h.GRID)
		for x := 0; x < h.GRID; x++ {
			t2[y][x] = make([]EntityPosition, len((*t)[y][x]))
			copy(t2[y][x], (*t)[y][x])

		}
	}
	return t2
}

// spawns a certain number of goroutines to iterate through entities, trying
// to lock them and get their position and send the entities and their
// positions to another set of goroutines handling the building of the
// list for each grid cell
func (h *SpatialHash) ComputeSpatialHash() {

	// this lock prevents another call to ComputeSpatialHash()
	// entering while we are currently calculating (this ensures robustness
	// if for some reason it is called too often)
	if !h.canEnterCompute.CAS(1, 0) {
		return
	}

	// set the computingTable pointer used by the receiver() workers to point
	// to the table we're building
	h.computingTable = &h.tables[h.nextBufIndex()]

	// prevent any updates to the spatialEntities list while we're using it
	h.spatialEntities.Mutex.Lock()
	defer h.spatialEntities.Mutex.Unlock()

	// acquire exclusive lock on the box component (position and bounding box)
	// TODO: should this happen at a higher level of abstraction, in the
	// system which will use the hash for physics / collision? do we really
	// want to let a bunch of position component lock acquires happen
	// in between computing the hash and using it? why not just lock the
	// position / hitbox components for the duration of the whole physics /
	// collision / spatial hash routine? (and lock velocity only for the
	// physics portion)
	// a similar lock will need to happen for the sprite component in Draw(),
	// but we can use the frozen position data from the spatial hash there
	h.em.Components.accessLocks[BOX_COMPONENT].Lock()
	defer h.em.Components.accessLocks[BOX_COMPONENT].Unlock()

	// set the number of entities remaining (used by receiver workers to
	// notify that computation is done)
	h.entitiesRemaining.Store(uint32(len(h.spatialEntities.Entities)))
	// clear the table data
	// NOTE: we "clear" the slice by setting its length to 0 (capacity remains
	// , so this is why a quadtree is a better structure if we'
	// re going to have entities clustering all into one place then
	// fanning out or clustering somewhere else
	for y := 0; y < h.GRID; y++ {
		for x := 0; x < h.GRID; x++ {
			cell := &((*h.computingTable)[y][x])
			*cell = (*cell)[:0]
		}
	}
	// Divide the list of entities into a certain number of partitions
	// which will be scanned by scanner() instances.
	// TODO: have the goroutines
	// already running, ready to be activated given a
	// range of entities to scan, each one pinned to a CPU
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
			go h.scanner(offset, partition_size)
		}
		<-h.computeDoneChannel
	}
	// if we're here, the computation has completed.
	// this increment, due to the modulo logic, is equivalent to setting
	// computedBufIndex = nextBufIndex
	h.tableBufIndex.Inc()
	// set canEnterCompute to 1 so the next call can enter and write
	h.canEnterCompute.Store(1)
}

// used to receive EntityPositions and put them into the right cell
func (h *SpatialHash) receiver(
	y int, x int,
	channel chan EntityPosition) {

	for {
		entityPosition := <-channel
		cell := &((*h.computingTable)[y][x])
		*cell = append(*cell, entityPosition)
		if h.entitiesRemaining.Dec() == 0 {
			h.computeDoneChannel <- true
		}
	}
}

// helper function for hashing
func (h *SpatialHash) cellForPoint(x int, y int) (cellX int, cellY int) {
	return x / (h.WORLD_HEIGHT / h.GRID), y / (h.WORLD_WIDTH / h.GRID)
}

// used to iterate the entities and send them to the right cell's channels
func (h *SpatialHash) scanner(offset int, partition_size int) {

	// iterate each entity in the partition
	for i := 0; i < partition_size; i++ {
		// get the entity's box
		entity := h.spatialEntities.Entities[offset+i]
		box := h.em.Components.Box[entity.ID]
		// find out how many grids the entity spans in x and y (almost always 0,
		// but we want to be thorough, and the fact that it's got a predictable
		// pattern 99% of the time means that branch prediction should help us)
		gridsHigh := int(box.H) / (h.WORLD_HEIGHT / h.GRID)
		gridsWide := int(box.W) / (h.WORLD_WIDTH / h.GRID)
		// figure out which cell the topleft corner is in
		topLeftCellX, topLeftCellY := h.cellForPoint(int(box.X), int(box.Y))
		// walk through each cell the entity touches by starting in the top-left
		// and walking according to gridsHigh and gridsWide
		for iy := 0; iy < gridsHigh+1; iy++ {
			for ix := 0; ix < gridsWide+1; ix++ {
				y := topLeftCellY + iy
				x := topLeftCellX + ix
				h.cellChannels[y][x] <- EntityPosition{entity, box}
			}
		}
	}
}

// turn a SpatialHashTable into a String representation (NOTE: do *NOT* call
// this on a pointer returned from CurrentTablePointer unless you can be sure
// that you have not called ComputeSpatialHash more than once - it does not
// lock the table it reads, and if you call ComputeSpatialHash twice, you may
// start to write to the table as this function reads it). Usually best to call
// on a Copy()
func (h *SpatialHash) String(table *SpatialHashTable) string {
	var buffer bytes.Buffer
	size := int(unsafe.Sizeof(*table))
	buffer.WriteString("[\n")
	for y := 0; y < h.GRID; y++ {
		for x := 0; x < h.GRID; x++ {
			cell := (*table)[y][x]
			size += int(unsafe.Sizeof(EntityPosition{})) * len(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				fmt.Sprintf("%+v", cell)))
			if !(y == h.GRID-1 && x == h.GRID-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}
