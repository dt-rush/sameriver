package engine

import (
	"fmt"
)

type DebugFunction func(s string, params ...interface{})

func genDebugFunction(moduleName string, flag bool) DebugFunction {
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

var entityManagerDebug = genDebugFunction(
	"EntityManager", DEBUG_ENTITY_MANAGER)
var eventDebug = genDebugFunction(
	"Events", DEBUG_EVENTS)
var updatedEntityListDebug = genDebugFunction(
	"UpdatedEntityList", DEBUG_UPDATED_ENTITY_LISTS)
var entityLogicDebug = genDebugFunction(
	"EntityLogic", DEBUG_ENTITY_LOGIC)
var goroutinesDebug = genDebugFunction(
	"Goroutines", DEBUG_GOROUTINES)
var atomicEntityModifyDebug = genDebugFunction(
	"AtomicModify", DEBUG_ATOMIC_MODIFY)
var entityClassDebug = genDebugFunction(
	"EntityClass", DEBUG_ENTITY_CLASS)
var worldLogicDebug = genDebugFunction(
	"WorldLogic", DEBUG_WORLD_LOGIC)
var entityLocksDebug = genDebugFunction(
	"EntityLocks", DEBUG_WORLD_LOGIC)
var despawnDebug = genDebugFunction(
	"Despawn", DEBUG_DESPAWN)
var behaviorDebug = genDebugFunction(
	"Behavior", DEBUG_BEHAVIOR)
