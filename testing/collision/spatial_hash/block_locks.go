package main

import (
	"sync"
)

func DoSpatialHash_Block(
	t *EntityTable,
	p *PositionComponent,
	buckets *SpatialHash) {

	// clear the buckes

	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			bucket := &buckets[y][x]
			*bucket = (*bucket)[:0]
		}
	}

	// calculate the spatial hash of each entity (use 12 goroutines), waiting
	// on locks
	partition_size := len(t.currentEntities) / GRID_GOROUTINES
	wg := sync.WaitGroup{}
	for part := 0; part < GRID_GOROUTINES; part++ {
		wg.Add(1)
		go func(part int) {
			offset := part * partition_size
			if part == GRID_GOROUTINES-1 {
				partition_size = len(t.currentEntities) - offset
			}
			for i := 0; i < partition_size; i++ {
				entity := t.currentEntities[offset+i]
				t.lockEntity(entity)
				position := p.Data[entity.ID]
				t.releaseEntity(entity)
				bucket := &buckets[position[1]/CELL_HEIGHT][position[0]/CELL_WIDTH]
				*bucket = append(*bucket, EntityPosition{entity, position})
			}
			wg.Done()
		}(part)
	}
	wg.Wait()
}
