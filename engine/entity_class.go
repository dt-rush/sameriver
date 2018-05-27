/*
 * EntityClass is used
 *
 */

package engine

import (
	"sync/atomic"
	"time"
)

type EntityClassDef struct {
	Name        string
	ListQueries []GenericEntityQuery
}

type EntityClass struct {
	Name  string
	Lists map[string](*UpdatedEntityList)
}

// Creates en EntityLogicFunc using a list of Behaviors
// See:
//		entity_logic_func.go (type definition)
//		logic_unit.go (containing type)
//		logic_component.go (storage for LogicUnits)
//		entity_manager.go (runs LogicUnits on entity spawn)
func (c *EntityClass) LogicUnitFromBehaviors(
	name string,
	behaviors []Behavior) LogicUnit {

	return NewLogicUnit(
		name,
		func(entity EntityToken,
			StopChannel chan bool,
			em *EntityManager) {

		logicloop:
			for {
				select {
				case <-StopChannel:
					break logicloop
				default:
					for i := 0; i < len(behaviors); i++ {
						if atomic.CompareAndSwapUint32(
							&(behaviors[i].running), 0, 1) {

							go func(behavior *Behavior) {
								behaviorDebug("Running behavior %s for entity "+
									"%d, ", behavior.Name, entity.ID)
								behavior.Func(entity, c, em)
								behaviorDebug("Sleeping %d ms for entity %d, "+
									"behavior: %s",
									behavior.Sleep.Nanoseconds()/1e6,
									entity.ID, behavior.Name)
								time.Sleep(behavior.Sleep)
								atomic.StoreUint32(&(behavior.running), 0)
							}(&behaviors[i])
						}
					}
				}
			}

		})
}
