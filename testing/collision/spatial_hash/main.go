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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var seconds = flag.Int("seconds", 10, "how many seconds to run for")
var strategy = flag.String("strategy", "",
	"which strategy to use for building the spatial hash (skip_locks,"+
		"block_locks, etc.)")

type Strategy func(t *EntityTable, p *PositionComponent, hash *SpatialHash)

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

	var F Strategy
	var FNAME string
	switch *strategy {
	case "skip_locks":
		computer := NewSpatialHashComputer()
		F = computer.DoSpatialHash_Skip
		FNAME = "DoSpatialHash_Skip"
	case "block_locks":
		F = DoSpatialHash_Block
		FNAME = "DoSpatialHash_Block"
	default:
		fmt.Println("not a valid spatial hash strategy! (needs one of " +
			"(skip_locks, block_locks, etc.)")
		os.Exit(1)
	}

	startTime := time.Now()
	rand.Seed(startTime.UnixNano())

	var Position PositionComponent
	var hash = NewSpatialHash()
	var entityTable EntityTable

	entityTable.SpawnEntities()
	DistributeEntities(&Position)
	StartBehaviors(&entityTable, &Position, *strategy)

	fmt.Println("testing spatial hashing...")
	fmt.Printf("time before loop started to run: %d ms\n",
		time.Since(startTime).Nanoseconds()/1e6)

	for i := 0; i < FRAMES_TO_RUN; i++ {
		if PRINT_EXAMPLE && i == FRAMES_TO_RUN/2 {
			fmt.Println(hash.String())
		}

		// preemptively prevent new locks
		time.Sleep(4 * time.Millisecond)
		entityTable.canLock.Store(0)
		t0 := time.Now()
		F(&entityTable, &Position, hash)
		milliseconds := time.Duration(time.Since(t0).Nanoseconds() / 1e6)
		entityTable.canLock.Store(1)
		if *cpuprofile == "" {
			fmt.Printf("%s took %d ms\n", FNAME, milliseconds)
		}
		if milliseconds > FRAME_SLEEP*2 {
			continue
		} else {
			time.Sleep(FRAME_SLEEP*2 - milliseconds)
		}
	}
}
