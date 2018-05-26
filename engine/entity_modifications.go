/*
 * Functions which can be passed to AtomicEntityModify or AtomicEntitiesModify
 *
 */

package engine

import ()

// user-facing despawn function which locks the EntityTable for a single
// despawn
func (m *EntityManager) Despawn(e EntityToken) {
	m.despawnInternal(e)
}

// a function which generates an entity modification to augment the health
// of an entity
func (m *EntityManager) ModifyHealth(change int) func(EntityToken) {
	return func(e EntityToken) {
		healthNow := int(m.Components.Health.Data[e.ID]) + change
		if healthNow > 255 {
			healthNow = 255
		} else if healthNow < 0 {
			healthNow = 0
		}
		m.Components.Health.Data[e.ID] = uint8(healthNow)
	}
}
