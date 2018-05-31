package engine

import (
	"fmt"
)

// produce a function which will act like fmt.Sprintf but be silent or not
// based on a supplied boolean value (below the function definition in this
// file you can find all of them used)
func DebugFunction(
	moduleName string, flag bool) func(s string, params ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", moduleName)
	return func(s string, params ...interface{}) {
		switch {
		case !flag:
			return
		case len(params) == 0:
			Logger.Printf(prefix + s)
		default:
			Logger.Printf(prefix+s, params...)
		}
	}
}

var entityManagerDebug = DebugFunction(
	"EntityManager", DEBUG_ENTITY_MANAGER)
var eventsDebug = DebugFunction(
	"Events", DEBUG_EVENTS)
var updatedEntityListDebug = DebugFunction(
	"UpdatedEntityList", DEBUG_UPDATED_ENTITY_LISTS)
var goroutinesDebug = DebugFunction(
	"Goroutines", DEBUG_GOROUTINES)
var atomicEntityModifyDebug = DebugFunction(
	"AtomicModify", DEBUG_ATOMIC_MODIFY)
var entityClassDebug = DebugFunction(
	"EntityClass", DEBUG_ENTITY_CLASS)
var worldLogicDebug = DebugFunction(
	"WorldLogic", DEBUG_WORLD_LOGIC)
var entityLocksDebug = DebugFunction(
	"EntityLocks", DEBUG_WORLD_LOGIC)
var spawnDebug = DebugFunction(
	"Spawn", DEBUG_SPAWN)
var despawnDebug = DebugFunction(
	"Despawn", DEBUG_DESPAWN)
var behaviorDebug = DebugFunction(
	"Behavior", DEBUG_BEHAVIOR)
var tagsDebug = DebugFunction(
	"Tags", DEBUG_TAGS)
var functionEndDebug = DebugFunction(
	">>>>>>>> Function End", DEBUG_FUNCTION_END)
var activeStateDebug = DebugFunction(
	"ActiveState", DEBUG_ACTIVE_STATE)
