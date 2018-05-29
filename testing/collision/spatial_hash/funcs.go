package main

import (
	"fmt"
	"github.com/dt-rush/donkeys-qquest/engine"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

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
	fmt.Printf("There are %d entities\n", len(t.currentEntities))
}

func (t *EntityTable) lockEntity(entity engine.EntityToken) {
	for t.canLock.Load() != 1 && !t.entityLocks[entity.ID].CAS(0, 1) {
		time.Sleep(FRAME_SLEEP / 2)
	}
}

func (t *EntityTable) lockEntity_blockable(entity engine.EntityToken) {
	for t.canLock.Load() != 1 && !t.entityLocks[entity.ID].CAS(0, 1) {
		time.Sleep(FRAME_SLEEP / 2)
	}
}

func (t *EntityTable) attemptLockEntity(entity engine.EntityToken) bool {
	return !t.entityLocks[entity.ID].CAS(0, 1)
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

func Behavior(t *EntityTable, entity engine.EntityToken, strategy string) {
	var LOCKFUNC func(engine.EntityToken)
	switch strategy {
	case "block_locks":
		LOCKFUNC = t.lockEntity_blockable
	default:
		LOCKFUNC = t.lockEntity
	}
	for {
		// simulating an AtomicEntityModify
		LOCKFUNC(entity)
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
			LOCKFUNC(otherEntity)
			time.Sleep(OTHER_ENTITY_LOCK_DURATION)
			t.releaseEntity(otherEntity)
		}
		t.releaseEntity(entity)
		time.Sleep(BEHAVIOR_POST_SLEEP())
	}
}

func StartBehaviors(t *EntityTable, p *PositionComponent, strategy string) {
	for i := 0; i < N_ENTITIES_WITH_BEHAVIOR; i++ {
		for j := 0; j < N_BEHAVIORS_PER_ENTITY; j++ {

			go Behavior(t, t.currentEntities[i%N_ENTITIES], strategy)
		}
	}
}
