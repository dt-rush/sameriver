package main

import (
	"fmt"
	"go.uber.org/atomic"
	"time"
)

// entity table functions as a kind of super lame drop-in replacement
// for the idea of the EntityManager from the engine
type EntityTable struct {
	entityLocks     [N_ENTITIES]atomic.Uint32
	currentEntities []EntityToken
}

func (t *EntityTable) SpawnEntities() {
	for i := 0; i < N_ENTITIES; i++ {
		t.currentEntities = append(t.currentEntities,
			EntityToken{ID: i})
	}
	fmt.Printf("There are %d entities\n", len(t.currentEntities))
}

func (t *EntityTable) lockEntity(entity EntityToken) {
	for !t.entityLocks[entity.ID].CAS(0, 1) {
		time.Sleep(FRAME_SLEEP / 2)
	}
}

func (t *EntityTable) attemptLockEntity(entity EntityToken) bool {
	return t.entityLocks[entity.ID].CAS(0, 1)
}

func (t *EntityTable) releaseEntity(entity EntityToken) {
	t.entityLocks[entity.ID].Store(0)
}
