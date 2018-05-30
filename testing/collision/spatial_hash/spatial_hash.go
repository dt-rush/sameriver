package main

import (
	"bytes"
	"fmt"
	"go.uber.org/atomic"
	"math"
	"time"
	"unsafe"
)

const GRID = 12
const SPATIAL_HASH_CELL_WIDTH = WORLD_WIDTH / GRID
const SPATIAL_HASH_CELL_HEIGHT = WORLD_HEIGHT / GRID

// used to store an entity with a position in a grid cell
type EntityPosition struct {
	entity   EntityToken
	position [2]int16
}

// the actual cell data structure is a GRID x GRID array of []EntityPosition
type SpatialHashTable [GRID][GRID][]EntityPosition

// deep-copy a spatial hash table
func (t *SpatialHashTable) Copy() SpatialHashTable {
	var t2 SpatialHashTable
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			t2[y][x] = make([]EntityPosition, len(t[y][x]))
			copy(t2[y][x], t[y][x])

		}
	}
	return t2
}

// turn a SpatialHashTable into a String representation (NOTE: do *NOT* call
// this on a pointer returned from CurrentTablePointer unless you can be sure
// that you have not called ComputeSpatialHash more than once - it does not
// lock the table it reads, and if you call ComputeSpatialHash twice, you may
// start to write to the table as this function reads it)
func (t *SpatialHashTable) String() string {
	var buffer bytes.Buffer
	size := int(unsafe.Sizeof(*t))
	buffer.WriteString("[\n")
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			cell := (*t)[y][x]
			size += int(unsafe.Sizeof(EntityPosition{})) * len(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				fmt.Sprintf("%+v", cell)))
			if !(y == GRID-1 && x == GRID-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}

// used to compute the spatial hash tables given a list of entities
type SpatialHash struct {
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
	// channels (one per grid square) used to receive successfully-read
	// positions from the goroutines which scan the entities
	// (not double-buffered since we exclude two
	// ComputeSpatialHash() instances from running at the same time using
	// the computeInProgress flag)
	cellChannels [GRID][GRID]chan EntityPosition
	// used to signal that the compute is done from one of the receive()
	// workers
	computeDoneChannel chan bool
	// an array of bools (one per buffer)
	// As we loop through partitions of the entity list, waiting to get all
	// entity positions to feed them to the goroutines appending to each
	// cell's list, we skip both those which are locked and, using this array,
	// those which have already been sent to the table-building goroutines
	alreadyRead [2][N_ENTITIES]bool // (N_ENTIITES = MAX_ENTITIES)
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
	// a reference to the position component
	position *PositionComponent
	// entityManager is used to acquire locks on entities
	entityTable *EntityTable
}

func NewSpatialHash(
	spatialEntities *UpdatedEntityList,
	entityTable *EntityTable,
	position *PositionComponent) *SpatialHash {

	h := SpatialHash{computeDoneChannel: make(chan bool)}
	// for each cell in the grid
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			// make a channel which will be used when building the grid
			h.cellChannels[y][x] = make(chan EntityPosition, N_ENTITIES)
			// prepare the expected slice storage for both SpatialHashTable
			// structures (we double-buffer) at this grid cell
			for buffer := 0; buffer < 2; buffer++ {
				// we assume a uniform distribution of entities
				// (N_ENTITIES / GRID^2 per cell)
				h.tables[buffer][y][x] = make([]EntityPosition,
					N_ENTITIES/(GRID*GRID))
			}
			// start the receiver for this cell
			go h.receiver(y, x, h.cellChannels[y][x])
		}
	}
	// take down references needed during compute
	h.spatialEntities = spatialEntities
	h.entityTable = entityTable
	h.position = position
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
	return h.CurrentTablePointer().Copy()
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

	// we don't want the UpdatedEntityList from modifying itself while
	// we read it
	h.spatialEntities.Mutex.Lock()
	defer h.spatialEntities.Mutex.Unlock()

	// set the number of entities remaining (used by receiver workers to
	// notify that computation is done)
	h.entitiesRemaining.Store(uint32(len(h.spatialEntities.Entities)))
	// clear the table data
	// NOTE: we "clear" the slice by setting its length to 0 (capacity remains
	// , so this is why a quadtree is a better structure if we'
	// re going to have entities clustering all into one place then
	// fanning out or clustering somewhere else
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			cell := &h.computingTable[y][x]
			*cell = (*cell)[:0]
		}
	}
	// Divide the list of entities into a certain number of partitions
	// which will be scanned by scanner() instances.
	// We determine the number of partitions via
	// N_PARTITIONS = 4 * log(N_ENTITIES+1)^2 + 1
	// We choose the number of partitions this way because it decently
	// approximates the estimated number of entities per cell assuming
	// a uniform distribution when the number of entities are in a reasonable
	// range, but also doesn't scale linearly with that number, approaching
	// a sort of soft "asymptote" around 50 for entity-counts less than
	// 4000 (that's a HELL OF A LOT OF ENTITIES!), or 75 entities
	// per partition. For 1600 entities that's 42 partitions with 38 entities
	// per partition.
	nScanPartitions := int(2*math.Pow(math.Log(N_ENTITIES+1), 2) + 1)
	partition_size := len(h.spatialEntities.Entities) / nScanPartitions
	for partition := 0; partition < nScanPartitions; partition++ {
		offset := partition * partition_size
		if partition == nScanPartitions-1 {
			// the last partition includes the remainder
			partition_size = len(h.spatialEntities.Entities) - offset
		}
		go h.scanner(offset, partition_size)
	}
	<-h.computeDoneChannel
	// if we're here, the computation has completed.
	// this increment, due to the modulo logic, is equivalent to setting
	// computedBufIndex = nextBufIndex
	h.tableBufIndex.Inc()
	// We will now spawn a goroutine to clear the `alreadyRead` data for the
	// next computation, and then set canEnterCompute to 1 so the next call
	// can enter and write
	go func() {
		for i := 0; i < N_ENTITIES; i++ {
			h.alreadyRead[h.nextBufIndex()][i] = false
			h.canEnterCompute.Store(1)
		}
	}()
}

// used to receive EntityPositions and put them into the right cell
func (h *SpatialHash) receiver(
	y int, x int,
	channel chan EntityPosition) {

	for {
		entityPosition := <-channel
		cell := &h.computingTable[y][x]
		*cell = append(*cell, entityPosition)
		if h.entitiesRemaining.Dec() == 0 {
			h.computeDoneChannel <- true
		}
	}
}

// used to iterate the entities and send them to the right cell's channels
func (h *SpatialHash) scanner(offset int, partition_size int) {

	// keep track of how many we've read
	n_read := 0
	for i := 0; n_read < partition_size; i = (i + 1) % partition_size {
		entity := h.spatialEntities.Entities[offset+i]
		if h.alreadyRead[h.nextBufIndex()][offset+i] {
			continue
		}
		// attempt the lock
		if h.entityTable.attemptLockEntity(entity) {
			// if we locked, grab the position and send it to
			// the channel
			position := h.position.Data[entity.ID]
			h.entityTable.releaseEntity(entity)
			y := position[1] / SPATIAL_HASH_CELL_HEIGHT
			x := position[0] / SPATIAL_HASH_CELL_WIDTH
			e := EntityPosition{entity, position}
			h.cellChannels[y][x] <- e
			h.alreadyRead[h.nextBufIndex()][offset+i] = true
			n_read++
			continue
		}
		// else, sleep a bit (to prevent hot loops if there are only
		// a few entities left and they are all locked)
		time.Sleep(10 * time.Microsecond)
	}
}
