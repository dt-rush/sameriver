package engine

import (
	"time"
)

// wait until we can lock an entity and check gen-match after lock. If
// gen mismatch, release and return false. else return true
func (m *EntityManager) lockEntity(entity EntityToken) bool {
	// wait until we can lock the entity
	for !m.entityTable.locks[entity.ID].CAS(0, 1) {
		// if we can't lock the entity, sleep half a frame
		time.Sleep(FRAME_SLEEP / 2)
	}
	// return value is whether we locked the entity with the same gen
	// (in between when the caller acquired the EntityToken and now, the
	// entity may have despawned)
	if !m.entityTable.genValidate(entity) {
		entityLocksDebug("LOCKED-GENMISMATCH entity %v", entity)
		m.releaseEntity(entity)
		return false
	}
	entityLocksDebug("LOCKED entity %v", entity)
	return true
}

// used by physics system to release an entity locked for modification
func (m *EntityManager) releaseEntity(entity EntityToken) {
	m.entityTable.locks[entity.ID].Store(0)
	entityLocksDebug("RELEASED entity %v", entity)
}

// lock multiple entities (with return value true only if gen matches for
// all entities locked)
func (m *EntityManager) lockEntities(entities []EntityToken) bool {
	// attempt to lock all entities, keeping track of which ones we have
	var allValid = true
	var locked = make([]EntityToken, 0)
	var time = time.Now().UnixNano()
	entityLocksDebug("[%d] attempting to lock %d entities: %v",
		time, len(entities), entities)
	for _, entity := range entities {
		if !m.lockEntity(entity) {
			entityLocksDebug("[%d] locking failed for %d", time, entity.ID)
			allValid = false
			break
		} else {
			entityLocksDebug("[%d] locking succeeded for %d", time, entity.ID)
			locked = append(locked, entity)
		}
	}
	// if one was invalid, the locking of this group no longer makes sense
	// (one was despawned since the []EntityToken was formulated by the caller)
	// so, release all those entities we've already locked and return false
	if !allValid {
		m.releaseEntities(locked)
		return false
	}
	// else return true
	return true
}

// release multiple entities
func (m *EntityManager) releaseEntities(entities []EntityToken) {
	for _, entity := range entities {
		m.releaseEntity(entity)
	}
}

// self explanatory
func (m *EntityManager) releaseTwoEntities(
	entityA EntityToken, entityB EntityToken) {
	m.entityTable.locks[entityA.ID].Store(0)
	m.entityTable.locks[entityB.ID].Store(0)
}

// used by physics system to attempt to lock an entity for modification, but
// will not sleep and retry if the lock fails, simply returns false if we
// didn't lock it, or if it had a new gen. Releases the entity if we locked it
// and gen doesn't match.
func (m *EntityManager) attemptLockEntityOnce(entity EntityToken) bool {
	// do a single attempt to lock
	locked := m.entityTable.locks[entity.ID].CAS(0, 1)
	// if we locked the entity but gen mistmatches, release it
	// and return false
	if locked &&
		!m.entityTable.genValidate(entity) {
		m.releaseEntity(entity)
		return false
	}
	// else, either we locked it and gen matched (pass), or we didn't lock
	// (fail), so return `locked`
	return locked
}

// used by collision system to attempt to lock two entities for modification
// (if both can't be acquired, we just back off and try again another cycle;
// collision between those two entities won't occur this cycle (there are many
// per second, so it's not noticeable to the user)
func (m *EntityManager) attemptLockTwoEntitiesOnce(
	entityA EntityToken, entityB EntityToken) bool {

	// attempt to lock entity A
	if !m.attemptLockEntityOnce(entityA) {
		// NOTE: we don't need to release the entity since if it failed
		// due to gen mismatch, attemptLockEntityOnce will itself release it,
		// and if it failed due to not locking, there's nothing to release
		return false
	}
	// attempt to lock entity B
	if !m.attemptLockEntityOnce(entityB) {
		// NOTE: if we're here, we *did* acquire A
		m.releaseEntity(entityA)
		return false
	}
	// if we're here, we locked both
	return true
}
