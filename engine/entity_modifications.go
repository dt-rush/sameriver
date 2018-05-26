/*
 * Functions which can be passed to AtomicEntityModify or AtomicEntitiesModify
 *
 */

package engine

import (
	"fmt"
	"time"
)

// user-facing despawn function which locks the EntityTable for a single
// despawn
func (m *EntityManager) Despawn(e EntityToken) {

	t0 := time.Now()
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()
	// check gen after mutex acquired, if it's different, this entity was
	// already despawned. mission accomplished! (this can happen if DespawnAll
	// is called, completes, and then we acquire the mutex)

	if DEBUG_DESPAWN {
		fmt.Printf("acquiring entityTable lock in despawn took: %d ms\n",
			time.Since(t0).Nanoseconds()/1e6)
	}
	m.despawnInternal(uint16(e.ID))
}
