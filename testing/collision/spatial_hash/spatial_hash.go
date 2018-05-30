package main

import (
	"bytes"
	"fmt"
	"go.uber.org/atomic"
	"sync"
	"time"
)

const GRID = 5
const CELL_WIDTH = WORLD_WIDTH / GRID
const CELL_HEIGHT = WORLD_HEIGHT / GRID
const SPATIAL_HASH_SCAN_GOROUTINES = 12
const SPATIAL_HASH_UPPER_ESTIMATE_ENTITIES_PER_SQUARE = N_ENTITIES / GRID
const GRID_GOROUTINES = 32

// used to store an entity with a position in a grid cell
type EntityPosition struct {
	entity   EntityToken
	position [2]int16
}

// used to compute the spatial hash cells given a list of entities
type SpatialHash struct {
	// a mutex to prevent from calculating a new hash before
	// the internal state has been reset
	resetMutex sync.Mutex
	// to keep track of whether we need to read an entity again as we
	// loop through the subsections waiting to get all entities in them
	alreadyRead [N_ENTITIES]bool // (N_ENTIITES = MAX_ENTITIES)
	// a channel used to receive successfully-read positions from the
	// goroutines this channel does not need to be allocated each time
	// we build
	cellChannels [GRID][GRID]chan EntityPosition
	// buckets is the actual spatial hash data structure
	cells [GRID][GRID][]EntityPosition
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	spatialEntities *UpdatedEntityList
	// a reference to the position component
	position *PositionComponent
	// entityManager is used to acquire locks on entities
	entityTable *EntityTable

	runningCellSenders   atomic.Uint32
	runningCellReceivers atomic.Uint32
}

func NewSpatialHash(
	spatialEntities *UpdatedEntityList,
	entityTable *EntityTable,
	position *PositionComponent) *SpatialHash {

	h := SpatialHash{}
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			h.cellChannels[y][x] = make(
				chan EntityPosition,
				N_ENTITIES)
			h.cells[y][x] = make(
				[]EntityPosition,
				SPATIAL_HASH_UPPER_ESTIMATE_ENTITIES_PER_SQUARE)
		}
	}
	h.spatialEntities = spatialEntities
	h.entityTable = entityTable
	h.position = position
	return &h
}

// resets the already read array (run after computing
func (h *SpatialHash) ResetAlreadyReadArray() {
	h.resetMutex.Lock()
	for i := 0; i < N_ENTITIES; i++ {
		h.alreadyRead[i] = false
	}
	h.resetMutex.Unlock()
}

// used to receive EntityPositions and put them into the right cell
func (h *SpatialHash) cellReceiver(
	y int, x int, entitiesRemaining *atomic.Uint32, wg *sync.WaitGroup) {

	for entitiesRemaining.Load() > 0 {
		select {
		case entityPosition := <-h.cellChannels[y][x]:
			h.cells[y][x] = append(h.cells[y][x], entityPosition)
			entitiesRemaining.Dec()
		default:
			time.Sleep(2 * time.Microsecond)
		}
	}
	wg.Done()
}

// used to iterate the entities and send them to the right cell's
// cellReceiver() instance (they are spawned, one for each cell, as goroutines)
func (h *SpatialHash) cellSender(
	offset int, partition_size int,
	wg *sync.WaitGroup) {

	// keep track of how many we've read
	n_read := 0
	for i := 0; n_read < partition_size; i = (i + 1) % partition_size {
		entity := h.spatialEntities.Entities[offset+i]
		if h.alreadyRead[offset+i] {
			continue
		}
		// attempt the lock
		if h.entityTable.attemptLockEntity(entity) {
			// if we locked, grab the position and send it to
			// the channel
			position := h.position.Data[entity.ID]
			h.entityTable.releaseEntity(entity)
			y := position[1] / CELL_HEIGHT
			x := position[0] / CELL_WIDTH
			h.cellChannels[y][x] <- EntityPosition{entity, position}
			h.alreadyRead[offset+i] = true
			n_read++
			continue
		}
		// else, sleep a bit (to prevent hot loops if there are only
		// a few entities left and they are all locked)
		time.Sleep(10 * time.Microsecond)
	}
	wg.Done()
}

// spawns a certain number of goroutines to iterate through entities, trying
// to lock them and get their position and send the entity to another goroutine
// handling the building of the list for that cell
func (h *SpatialHash) ComputeSpatialHash() {

	// lock the UpdatedEntityList from modifying itself while we read it
	h.spatialEntities.Mutex.Lock()
	defer h.spatialEntities.Mutex.Unlock()

	// mutex / deferred goroutine pattern means that
	// the reset of internal state used to check entities
	// will be done immediately after the function returns,
	// and in the background
	h.resetMutex.Lock()
	h.resetMutex.Unlock()
	defer func() {
		go h.ResetAlreadyReadArray()
	}()

	// waitgroup is used to ensure that every entity has been checked and every
	// grid has bene populated
	wg := sync.WaitGroup{}
	// divide the list of entities into a certain number of partitions
	partition_size := len(h.spatialEntities.Entities) /
		SPATIAL_HASH_SCAN_GOROUTINES
	for partition := 0; partition < SPATIAL_HASH_SCAN_GOROUTINES; partition++ {
		// the last partition includes the remainder
		offset := partition * partition_size
		if partition == SPATIAL_HASH_SCAN_GOROUTINES-1 {
			partition_size = len(h.spatialEntities.Entities) - offset
		}
		// spawn the cellSender goroutine for this partition
		wg.Add(1)
		go h.cellSender(offset, partition_size, &wg)
	}
	// for each cell, spawn a cellReceiver goroutine
	entitiesRemaining := atomic.NewUint32(N_ENTITIES)
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			wg.Add(1)
			go h.cellReceiver(y, x, entitiesRemaining, &wg)
		}
	}
	wg.Wait()
}

func (h *SpatialHash) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %v", x, y, h.cells[y][x]))
			if !(y == GRID-1 && x == GRID-1) {
				buffer.WriteString(",")
			}
			if !(y == GRID-1) {
				buffer.WriteString("\n")
			}
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}
