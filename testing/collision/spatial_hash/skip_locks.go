package main

import (
	"fmt"
	"sync"
	"time"
)

type SpatialHashComputer struct {
	// a mutex to prevent from calculating a new hash before
	// the internal state has been reset
	resetMutex sync.Mutex
	// to keep track of whether we need to read an entity again as we
	// loop through the subsections waiting to get all entities in them
	alreadyRead []bool
	// a channel used to receive successfully-read positions from the
	// goroutines this channel does not need to be allocated each time
	// we build
	builderChannel chan EntityPosition
}

func NewSpatialHashComputer() *SpatialHashComputer {
	fmt.Println("building computer state for DoSpatialHash_Skip")
	c := SpatialHashComputer{
		sync.Mutex{},
		make([]bool, N_ENTITIES),              // MAX_ENTITIES
		make(chan EntityPosition, N_ENTITIES), // MAX_ENTITIES
	}
	return &c
}

func (c *SpatialHashComputer) ResetInternalState() {
	c.resetMutex.Lock()
	for i := 0; i < N_ENTITIES; i++ {
		c.alreadyRead[i] = false
	}
	c.resetMutex.Unlock()
}

func (c *SpatialHashComputer) DoSpatialHash_Skip(
	t *EntityTable,
	p *PositionComponent,
	buckets *SpatialHash) {

	fmt.Println("DoSpatialHash_Skip")

	c.resetMutex.Lock()
	c.resetMutex.Unlock()

	// when the function ends, spawn a goroutine to clean its internal
	// state
	defer func() {
		go c.ResetInternalState()
	}()

	// clear the data structure
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			bucket := &buckets[y][x]
			*bucket = (*bucket)[:0]
		}
	}

	// calculate the spatial hash of each entity (use 12 goroutines)
	// Skip locked entities

	n_entities := len(t.currentEntities)
	partitions := 12
	partition_size := n_entities / partitions
	// wait group used to signal that we have built the data structure
	wg := sync.WaitGroup{}
	for part := 0; part < partitions; part++ {
		// start one subsection's goroutine
		wg.Add(1)
		go func(part int) {
			// size the last partition properly
			offset := part * partition_size
			if part == GRID_GOROUTINES-1 {
				partition_size = n_entities - offset
			}
			// keep track of how many we've read
			n_read := 0
			i := 0
			fmt.Printf("part %d of %d is [%d, %d)\n",
				part+1, partitions, offset, offset+partition_size)
			for n_read < partition_size {
				entity := t.currentEntities[offset+i]
				quickskip := c.alreadyRead[offset+i]
				if quickskip {
					i = (i + 1) % partition_size
					continue
				}
				// attempt the lock
				if t.attemptLockEntity(entity) {
					// if we locked, grab the position and send it to
					// the channel
					position := p.Data[entity.ID]
					t.releaseEntity(entity)
					c.builderChannel <- EntityPosition{entity, position}
					c.alreadyRead[offset+i] = true
					n_read++
					i = (i + 1) % partition_size
					continue
				}
				// else, sleep a bit (to prevent hot loops if there are only
				// a few entities left and they are all locked)
				time.Sleep(10 * time.Microsecond)
				i = (i + 1) % partition_size
				continue
			}
			fmt.Printf("part %d completed.\n", part+1)
			wg.Done()
		}(part)
	}
	// create a goroutine to read the entities as they come out the channel
	wg.Add(1)
	go func() {
		n_read := 0
		for n_read < n_entities {
			ep := <-c.builderChannel
			bucket := &buckets[ep.position[1]/CELL_HEIGHT][ep.position[0]/CELL_WIDTH]
			*bucket = append(*bucket, ep)
			n_read++
		}
		wg.Done()
	}()
	wg.Wait()
}
