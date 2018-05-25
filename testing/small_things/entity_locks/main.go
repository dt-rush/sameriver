package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const N_ENTITIES = 1024
const N_LOCKERS = 1024
const SLEEP_MS = 25
const N_LOOPS = 32

type Uint32LockTable struct {
	entityLocks [N_ENTITIES]uint32
}

func locker(table *Uint32LockTable, cond *sync.Cond, wg *sync.WaitGroup) {

	stopChan := make(chan bool)

	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		stopChan <- true
	}()

loop:
	for {
		select {
		case _ = <-stopChan:
			break loop
		default:
			go func() {
				id := rand.Intn(N_ENTITIES)
				// wait for lock
				for !atomic.CompareAndSwapUint32(
					&table.entityLocks[id], 0, 1) {

					time.Sleep(time.Duration(
						rand.Intn(2*SLEEP_MS)) * time.Millisecond)
				}
				// sleep
				time.Sleep(time.Duration(
					rand.Intn(2*SLEEP_MS)) * time.Millisecond)
				// unset lock
				atomic.StoreUint32(&table.entityLocks[id], 0)
			}()
			time.Sleep(time.Duration(
				rand.Intn(2*SLEEP_MS)) * time.Millisecond)
		}
	}

	wg.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)
	table := Uint32LockTable{}
	var counters [N_ENTITIES]uint32
	fmt.Printf("%d entities, %d lockers, sleep %d ms\n",
		N_ENTITIES, N_LOCKERS, SLEEP_MS)
	fmt.Printf("will check each entity in a loop %d times\n",
		N_LOOPS)
	for i := 0; i < N_LOCKERS; i++ {
		wg.Add(1)
		go locker(&table, cond, &wg)
	}
	sumMilliseconds := float64(0)
	for i := 0; i < N_LOOPS; i++ {
		t0 := time.Now().UnixNano()
		for j := 0; j < N_ENTITIES; j++ {
			if atomic.CompareAndSwapUint32(&table.entityLocks[j], 0, 0) {
				if counters[j] == 0 {
					counters[j]++
				} else {
					counters[j] *= 2
				}
			}
		}
		t1 := time.Now().UnixNano()
		elapsed := float64(t1-t0) / float64(1e6)
		sumMilliseconds += elapsed
		time.Sleep(16 * time.Millisecond)
	}
	cond.Broadcast()
	wg.Wait()
	fmt.Printf("average time to check %d Uint32's with "+
		"atomic.CompareAndSwapUint32: %.3f ms\n",
		N_ENTITIES,
		sumMilliseconds/float64(N_LOOPS))
}
