/*
 * Behavior represents a certain logic which will be run for an entity
 * with a sleep time between each run. The entity is inherently a member
 * of an entity class, even if it's a singleton.
 *
 */

package engine

import (
	"go.uber.org/atomic"
	"time"
)

// the type of a function run
type BehaviorFunc func(
	e EntityToken,
	em *EntityManager)

type Behavior struct {
	Name string
	// a constant amount of time to sleep after each time Func is run
	Sleep time.Duration
	// the function this behaviour represents (run when running is 0)
	Func BehaviorFunc
	// used atomically as a lock to determine whether to run the Func
	running atomic.Uint32
}

// Creates en EntityLogicFunc using a list of Behaviors
// See:
//		entity_logic_func.go (type definition)
//		logic_unit.go (containing type)
//		logic_component.go (storage for LogicUnits)
//		entity_manager.go (runs LogicUnits on entity spawn)
func LogicUnitFromBehaviors(
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
						if behaviors[i].running.CAS(0, 1) {

							go func(behavior *Behavior) {
								behaviorDebug("Running behavior %s for entity "+
									"%d, ", behavior.Name, entity.ID)
								behavior.Func(entity, em)
								behaviorDebug("Sleeping %d ms for entity %d, "+
									"behavior: %s",
									behavior.Sleep.Nanoseconds()/1e6,
									entity.ID, behavior.Name)
								time.Sleep(behavior.Sleep)
								behavior.running.Store(0)
							}(&behaviors[i])
						}
					}
					// we need to sleep here in order to avoid burning the CPU!
					time.Sleep(5 * FRAME_SLEEP)
				}
			}

		})
}