package engine

import (
	"fmt"
	"sync/atomic"
	"time"
)

type EntityClass struct {
	Lists map[string](*UpdatedEntityList)
	Name  string
}

func NewEntityClass(
	em *EntityManager,
	className string,
	queryLists []GenericEntityQuery) EntityClass {

	entityClassDebug("in NewEntityClass(%s)", className)
	Lists := make(map[string](*UpdatedEntityList), len(queryLists))
	for _, q := range queryLists {
		entityClassDebug("trying to build list %s for class %s",
			q.Name, className)
		Lists[q.Name] = em.GetUpdatedActiveEntityList(
			q, fmt.Sprintf("class(%s):%s", className, q.Name))
	}
	return EntityClass{
		Name:  className,
		Lists: Lists}
}

type BehaviorFunc func(
	e EntityToken,
	c *EntityClass,
	em *EntityManager)

type Behavior struct {
	// a constant amount of time to sleep after each time Func is run
	Sleep time.Duration
	// the function this behaviour represents (run when running is 0)
	Func BehaviorFunc
	// used atomically as a lock to determine whether to run the Func
	running uint32
}

func (c *EntityClass) GenerateLogicFunc(
	behaviors []Behavior) LogicFunc {

	return func(entity EntityToken,
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
							behavior.Func(entity, c, em)
							entityClassDebug("Sleeping %d ms for entity %d, "+
								"behavior %d",
								behavior.Sleep.Nanoseconds()/1e6, entity.ID, i)
							time.Sleep(behavior.Sleep)
							atomic.StoreUint32(&(behavior.running), 0)
						}(&behaviors[i])
					}
				}
			}
		}
	}
}
