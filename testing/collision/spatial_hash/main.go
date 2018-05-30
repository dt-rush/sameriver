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

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var FRAMES_TO_RUN = FPS * *seconds

	startTime := time.Now()
	rand.Seed(startTime.UnixNano())

	// build the position component
	var Position PositionComponent
	for i := 0; i < N_ENTITIES; i++ {
		Position.Data[i] = [2]int16{
			int16(rand.Intn(WORLD_WIDTH)),
			int16(rand.Intn(WORLD_HEIGHT))}
	}
	// build an entity table
	var entityTable EntityTable
	entityTable.SpawnEntities()
	// build a list of spatial entities
	spatialEntities := UpdatedEntityList{Entities: entityTable.currentEntities}
	fmt.Printf("entityTable.currentEntities: %d\n",
		len(entityTable.currentEntities))
	// build the spatial hash computer
	spatialHash := NewSpatialHash(&spatialEntities, &entityTable, &Position)
	// start the behavior goroutines
	StartBehaviors(&entityTable, &Position)

	fmt.Println("testing spatial hashing...")
	fmt.Printf("time before loop started to run: %d ms\n",
		time.Since(startTime).Nanoseconds()/1e6)

	for i := 0; i < FRAMES_TO_RUN; i++ {
		t0 := time.Now()
		if PRINT_EXAMPLE && i == FRAMES_TO_RUN/2 {
			fmt.Println(spatialHash.String())
		}
		// preemptively prevent new locks
		time.Sleep(2 * time.Millisecond)
		t1 := time.Now()
		spatialHash.ComputeSpatialHash()
		milliseconds := time.Duration(time.Since(t1).Nanoseconds() / 1e6)
		if *cpuprofile == "" {
			fmt.Printf("computing spatial hash took %d ms\n", milliseconds)
		}
		if *printHash {
			fmt.Println(spatialHash.String())
		}
		totalMilliseconds := time.Duration(time.Since(t0).Nanoseconds() / 1e6)
		if totalMilliseconds >= FRAME_SLEEP {
			continue
		} else {
			time.Sleep(FRAME_SLEEP - totalMilliseconds)
		}
	}
}
