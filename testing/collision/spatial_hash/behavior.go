package main

import (
	"math/rand"
	"time"
)

// imitates a behavior func which will atomically modify the entity,
// possibly run a safeget on another entity, possibly lock that other
// entity, and sleep various amounts of time during all this
func Behavior(t *EntityTable, entity EntityToken) {
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
			otherEntity := EntityToken{ID: otherID}
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
