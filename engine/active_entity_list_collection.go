package engine


type ActiveEntityListCollection struct {
	em *EntityManager
	watchers []*QueryWatcher
	lists []*UpdatedEntityList
	mutex sync.RWMutex
}

func (c *ActiveEntityListCollection) Init(em *EntityManager) {
	c.em = em
}

func (c *ActiveEntityListCollection) CreateUpdatedEntityList(q EntityQuery) {
	// register a query watcher for the query given
	qw := NewEntityQueryWatcher(q)
	m.watchers = append(m.watchers, qw)
	// build the list (as yet empty), provide it with a backlog, and start it
	backlog := c.em.entityTable.snapshotAllocatedEntities()
	backlogTester := func(entity EntityToken) bool {
		return q.Test(entity, c.em)
	}
	list := NewUpdatedEntityList(qw, backlog, backlogTester)
	list.start()
	return list
}

func (c *ActiveEntityListCollection) notifyActiveState(
	entity EntityToken, active bool) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, watcher := range c.watchers {
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
					"remove(deactivated):%v to %s", time, entity, watcher.Name)
				watcher.Channel <- EntitySignal{ENTITY_REMOVE, entity}
			} else {
				updatedEntityListDebug("sending "+
					"add(activated):%v to %s", time, entity, watcher.Name)
				watcher.Channel <- EntitySignal{ENTITY_ADD, entity}
			}
		}
	}


}

