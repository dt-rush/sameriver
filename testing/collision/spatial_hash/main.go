package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func PrintIfNotProfiling(s string, args ...interface{}) {
	if *cpuprofile == "" {
		fmt.Printf(s, args...)
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// parse CLI flags and set up profiling if need be
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// the CLI flag "-seconds" is the number of seconds to run for
	var FRAMES_TO_RUN = FPS * *seconds

	// build the position component
	var Position PositionComponent
	for i := 0; i < N_ENTITIES; i++ {
		Position.Data[i] = [2]int16{
			int16(rand.Intn(WORLD_WIDTH)),
			int16(rand.Intn(WORLD_HEIGHT))}
	}
	// prepare data needed for simulation
	var entityTable EntityTable
	entityTable.SpawnEntities()
	spatialEntities := UpdatedEntityList{Entities: entityTable.currentEntities}
	PrintIfNotProfiling("entityTable.currentEntities: %d\n",
		len(entityTable.currentEntities))
	spatialHash := NewSpatialHash(&spatialEntities, &entityTable, &Position)
	// start a bunch of goroutines to lock entities
	StartBehaviors(&entityTable, &Position)

	PrintIfNotProfiling("testing spatial hashing...")

	for i := 0; i < FRAMES_TO_RUN; i++ {
		// do the compute
		t0 := time.Now()
		spatialHash.ComputeSpatialHash()
		computeMilliseconds := time.Duration(time.Since(t0).Nanoseconds() / 1e6)
		PrintIfNotProfiling("computing spatial hash took %d ms\n",
			computeMilliseconds)
		PrintIfNotProfiling("table pointer: %p\n", spatialHash.CurrentTable())
		// do a copy
		t1 := time.Now()
		currentTable := spatialHash.CurrentTableCopy()
		copyMilliseconds := time.Duration(time.Since(t1).Nanoseconds() / 1e6)
		PrintIfNotProfiling("copying spatial hash table took %d ms\n",
			copyMilliseconds)
		// print the hash halfway through, if -printhash passed
		if *printHash && i == FRAMES_TO_RUN/2 {
			PrintIfNotProfiling("%s\n", currentTable.String())
		}
		// determine how long to sleep for in order to make good on 16 ms per
		// loop
		totalMilliseconds := time.Duration(time.Since(t0).Nanoseconds() / 1e6)
		if totalMilliseconds >= FRAME_SLEEP {
			continue
		} else {
			time.Sleep(FRAME_SLEEP - totalMilliseconds)
		}
	}
}
