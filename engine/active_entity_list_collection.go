package engine

import (
	"sync"
)

type ActiveEntityListCollection struct {
	em       *EntityManager
	watchers map[string]*EntityQueryWatcher
	lists    map[string]*UpdatedEntityList
	mutex    sync.RWMutex
}

func (c *ActiveEntityListCollection) Init(em *EntityManager) {
	c.em = em
	c.watchers = make(map[string]*EntityQueryWatcher)
	c.lists = make(map[string]*UpdatedEntityList)
}

func (c *ActiveEntityListCollection) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {

	updatedEntityListDebug("waiting for ActiveEntityListCollection mutex "+
		"in GetUpdatedEntityList for %s...", q.Name)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	updatedEntityListDebug("acquired ActiveEntityListCollection mutex "+
		"in GetUpdatedEntityList for %s...", q.Name)

	// return the list if it already exists (this is why query names should
	// be unique if they expect to be unique!)
	// TODO: document this requirement
	if list, exists := c.lists[q.Name]; exists {
		updatedEntityListDebug("list for query named \"%s\" already "+
			"exists. Returning that list.", q.Name)
		return list
	}

	// register a query watcher for the query given
	qw := NewEntityQueryWatcher(q)
	c.watchers[q.Name] = &qw
	// build the list (as yet empty), provide it with a backlog, and start it
	backlog := c.em.entityTable.snapshotAllocatedEntities()
	backlogTester := func(entity EntityToken) bool {
		return q.Test(entity, c.em)
	}
	list := NewUpdatedEntityList(qw, backlog, backlogTester)
	list.start()
	c.lists[q.Name] = list
	updatedEntityListDebug("returning list %s, which will now start to build",
		q.Name)
	return list
}

func (c *ActiveEntityListCollection) notifyActiveState(
	entity EntityToken, active bool) {

	updatedEntityListDebug("waiting for mutex in notifyActiveState for %v...", entity)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer updatedEntityListDebug("released mutex in notifyActiveState for %v", entity)

	updatedEntityListDebug("acquired mutex in notifyActiveState for %v", entity)

	for _, watcher := range c.watchers {

		go func(watcher *EntityQueryWatcher) {
			updatedEntityListDebug("in notifyActiveState goroutine, "+
				"testing query %s for %v...", watcher.Query.Name, entity)
			if watcher.Query.Test(entity, c.em) {
				// warn if the channel is full (we will block here if so)
				// NOTE: this can be very bad indeed, since now whatever
				// called Activate is blocking
				if len(watcher.Channel) == ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY {
					entityManagerDebug("⚠  active entity "+
						" watcher channel %s is full, causing block in "+
						" for NotifyActiveState(%d, %v)\n",
						watcher.Name, entity.ID, active)
				}
				// send the signal
				if !active {
					updatedEntityListDebug("sending "+
						"remove(deactivated):%v to %s", entity, watcher.Name)
					watcher.Channel <- EntitySignal{ENTITY_REMOVE, entity}
				} else {
					updatedEntityListDebug("sending "+
						"add(activated):%v to %s", entity, watcher.Name)
					watcher.Channel <- EntitySignal{ENTITY_ADD, entity}
				}
			}
		}(watcher)
	}
}

func (c *ActiveEntityListCollection) checkActiveEntity(entity EntityToken) {

	updatedEntityListDebug("waiting for mutex in checkActiveEntity for %v...", entity)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer updatedEntityListDebug("released mutex in checkActiveEntity for %v", entity)

	updatedEntityListDebug("acquired mutex in checkActiveEntity for %v", entity)

	for _, watcher := range c.watchers {

		go func(watcher *EntityQueryWatcher) {
			updatedEntityListDebug("in checkActiveEntity goroutine, "+
				"testing query %s for %v...", watcher.Query.Name, entity)
			if watcher.Query.Test(entity, c.em) {
				// warn if the channel is full (we will block here if so)
				// NOTE: this can be very bad indeed, since now whatever
				// called Activate is blocking
				if len(watcher.Channel) == ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY {
					entityManagerDebug("⚠  active entity "+
						" watcher channel %s is full, causing block in "+
						" for checkActiveEntity(%v)\n",
						watcher.Name, entity.ID)
				}
				// send the signal
				updatedEntityListDebug("sending "+
					"add(checked):%v to %s", entity, watcher.Name)
				watcher.Channel <- EntitySignal{ENTITY_ADD, entity}
			}
		}(watcher)
	}
}
