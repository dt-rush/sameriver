package engine

type ActiveEntityListCollection struct {
	em    *EntityManager
	lists map[string]*UpdatedEntityList
}

func NewActiveEntityListCollection(
	em *EntityManager) *ActiveEntityListCollection {

	return &ActiveEntityListCollection{
		em:    em,
		lists: make(map[string]*UpdatedEntityList),
	}
}

func (c *ActiveEntityListCollection) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {
	return c.getUpdatedEntityList(q, false)
}

func (c *ActiveEntityListCollection) GetSortedUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {
	return c.getUpdatedEntityList(q, true)
}

func (c *ActiveEntityListCollection) getUpdatedEntityList(
	q EntityQuery, sorted bool) *UpdatedEntityList {
	// return the list if it already exists (this is why query names should
	// be unique if they expect to be unique!)
	// TODO: document this requirement
	if list, exists := c.lists[q.Name]; exists {
		return list
	}
	// register a query watcher for the query given
	var list *UpdatedEntityList
	if sorted {
		list = NewSortedUpdatedEntityList()
	} else {
		list = NewUpdatedEntityList()
	}
	list.Query = &q
	c.processBacklog(q, list)
	c.lists[q.Name] = list
	return list
}

func (c *ActiveEntityListCollection) processBacklog(
	q EntityQuery,
	list *UpdatedEntityList) {

	for _, e := range c.em.GetCurrentEntities() {
		if q.Test(e, c.em) {
			list.Signal(EntitySignal{ENTITY_ADD, e})
		}
	}
}

func (c *ActiveEntityListCollection) notifyActiveState(
	entity *EntityToken, active bool) {

	// send add / remove signal to all lists
	for _, list := range c.lists {
		if list.Query.Test(entity, c.em) {
			if active {
				list.Signal(EntitySignal{ENTITY_ADD, entity})
			} else {
				list.Signal(EntitySignal{ENTITY_REMOVE, entity})
			}
		}
	}
}

func (c *ActiveEntityListCollection) checkActiveEntity(entity *EntityToken) {

	// check if the entity needs to be added to any lists
	for _, list := range c.lists {
		if list.Query.Test(entity, c.em) {
			list.Signal(EntitySignal{ENTITY_ADD, entity})
		}
	}
	// check whether the entity needs to be removed from any lists it's on
	toRemove := make([]*UpdatedEntityList, 0)
	for _, list := range entity.Lists {
		if list.Query != nil && !list.Query.Test(entity, c.em) {
			toRemove = append(toRemove, list)
		}
	}
	for _, list := range toRemove {
		list.Signal(EntitySignal{ENTITY_REMOVE, entity})
	}
}
