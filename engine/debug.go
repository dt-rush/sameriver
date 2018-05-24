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
	"EventManager", DEBUG_EVENTS)
var updatedEntityListDebug = genDebugFunction(
	"UpdatedEntityList", DEBUG_UPDATED_ENTITY_LISTS)
var logicDebug = genDebugFunction(
	"Logic", DEBUG_LOGIC)
var goroutinesDebug = genDebugFunction(
	"Logic", DEBUG_GOROUTINES)
var atomicEntityModifyDebug = genDebugFunction(
	"AtomicModify", DEBUG_ATOMIC_MODIFY)
