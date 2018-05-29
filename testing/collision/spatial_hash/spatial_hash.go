package main

import (
	"flag"
	"fmt"
	"github.com/dt-rush/donkeys-qquest/engine"
	"go.uber.org/atomic"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

const N_ENTITIES = 1024
const N_ENTITIES_WITH_BEHAVIOR = 800
const N_BEHAVIORS_PER_ENTITY = 4
const CHANCE_TO_SAFEGET_IN_BEHAVIOR = 0.3
const SAFEGET_DURATION = 20 * time.Microsecond
const CHANCE_TO_LOCK_OTHER_ENTITY = 0.05
const OTHER_ENTITY_LOCK_DURATION = 20 * time.Microsecond

const WORLD_WIDTH = 3200
const WORLD_HEIGHT = 3200
const GRID = 32
const CELL_WIDTH = WORLD_WIDTH / GRID
const CELL_HEIGHT = WORLD_HEIGHT / GRID
const UPPER_ESTIMATE_ENTITIES_PER_SQUARE = 24

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond
const FRAMES_TO_RUN = FPS * 10

var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(4000)) * time.Nanosecond
}

var BEHAVIOR_POST_SLEEP = func() time.Duration {
	return time.Duration(time.Duration(rand.Intn(700)) * time.Millisecond)
}

const PRINT_EXAMPLE = false

type PositionComponent struct {
	Data [N_ENTITIES][2]int16
}

type EntityTable struct {
	canLock         atomic.Uint32
	entityLocks     [N_ENTITIES]atomic.Uint32
	currentEntities []engine.EntityToken
}

func (t *EntityTable) SpawnEntities() {
	for i := 0; i < N_ENTITIES; i++ {
		t.currentEntities = append(t.currentEntities,
			engine.EntityToken{ID: i})
	}
}

func (t *EntityTable) lockEntity(entity engine.EntityToken) {
	for t.canLock.Load() != 1 && !t.entityLocks[entity.ID].CAS(0, 1) {
		time.Sleep(FRAME_SLEEP / 2)
	}
}

func (t *EntityTable) releaseEntity(entity engine.EntityToken) {
	t.entityLocks[entity.ID].Store(0)
}

func DistributeEntities(p *PositionComponent) {
	for i := 0; i < N_ENTITIES; i++ {
		p.Data[i] = [2]int16{
			int16(rand.Intn(WORLD_WIDTH)),
			int16(rand.Intn(WORLD_HEIGHT))}
	}
}

func Behavior(t *EntityTable, entity engine.EntityToken) {
	for {
		// simulating an AtomicEntityModify
		t.lockEntity(entity)
		time.Sleep(ATOMIC_MODIFY_DURATION())
		if rand.Float32() < CHANCE_TO_SAFEGET_IN_BEHAVIOR {
			time.Sleep(SAFEGET_DURATION)
		}
		if rand.Float32() < CHANCE_TO_LOCK_OTHER_ENTITY {
			otherID := rand.Intn(N_ENTITIES)
			for otherID == entity.ID {
				otherID = rand.Intn(N_ENTITIES)
			}
			otherEntity := engine.EntityToken{ID: otherID}
			t.lockEntity(otherEntity)
			time.Sleep(OTHER_ENTITY_LOCK_DURATION)
			t.releaseEntity(otherEntity)
		}
		t.releaseEntity(entity)
		time.Sleep(BEHAVIOR_POST_SLEEP())
	}
}

func StartBehaviors(t *EntityTable, p *PositionComponent) {
	for i := 0; i < N_ENTITIES_WITH_BEHAVIOR; i++ {
		for j := 0; j < N_BEHAVIORS_PER_ENTITY; j++ {
			go Behavior(t, t.currentEntities[i%N_ENTITIES])
		}
	}
}

func AllocateBuckets(buckets *[GRID][GRID][]engine.EntityToken) {
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			buckets[y][x] = make([]engine.EntityToken,
				UPPER_ESTIMATE_ENTITIES_PER_SQUARE)
		}
	}
}

func DoSpatialHash(
	t *EntityTable,
	p *PositionComponent,
	buckets *[GRID][GRID][]engine.EntityToken) {

	// clear the buckes

	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			bucket := &buckets[y][x]
			*bucket = (*bucket)[:0]
		}
	}

	// calculate the spatial hash of each entity (use 12 goroutines
	// because getting locked on an entity sucks)
	// TODO: could improve by skipping entities we can't lock and returning
	// to them once we've processed all the others
	partitions := 12
	partition_size := len(t.currentEntities) / partitions
	wg := sync.WaitGroup{}
	for part := 0; part < partitions; part++ {
		wg.Add(1)
		go func(part int) {
			offset := part * partition_size
			for i := 0; i < partition_size; i++ {
				entity := t.currentEntities[offset+i]
				t.lockEntity(entity)
				position := p.Data[entity.ID]
				t.releaseEntity(entity)
				bucket := &buckets[position[1]/CELL_HEIGHT][position[0]/CELL_WIDTH]
				*bucket = append(*bucket, entity)
			}
			wg.Done()
		}(part)
	}
	wg.Wait()
}

func PrintBuckets(buckets *[GRID][GRID][]engine.EntityToken) {
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			fmt.Printf("GRID(%d, %d): %v\n\n", x, y, buckets[y][x])
		}
	}
}

func main() {

	startTime := time.Now()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Now().UnixNano())

	var Position PositionComponent
	var buckets [GRID][GRID][]engine.EntityToken
	var entityTable EntityTable

	AllocateBuckets(&buckets)
	entityTable.SpawnEntities()
	DistributeEntities(&Position)
	StartBehaviors(&entityTable, &Position)

	fmt.Println("testing spatial hashing...")
	fmt.Printf("time before loop started to run: %d ms\n",
		time.Since(startTime).Nanoseconds()/1e6)

	for i := 0; i < FRAMES_TO_RUN; i++ {
		if PRINT_EXAMPLE && i == FRAMES_TO_RUN/2 {
			PrintBuckets(&buckets)
		}

		entityTable.canLock.Store(0)
		time.Sleep(4 * time.Millisecond)
		t0 := time.Now()
		DoSpatialHash(&entityTable, &Position, &buckets)
		milliseconds := time.Duration(time.Since(t0).Nanoseconds() / 1e6)
		entityTable.canLock.Store(1)
		// fmt.Printf("took %d ms to do spatial hash\n", milliseconds)
		if milliseconds > FRAME_SLEEP*2 {
			continue
		} else {
			time.Sleep(FRAME_SLEEP*2 - milliseconds)
		}
	}
}
