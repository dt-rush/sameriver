package engine

type EntityQuery interface {
	Test(id uint16, em *EntityManager) bool
}

type GenericEntityQuery struct {
	Name     string
	TestFunc func(id uint16, em *EntityManager) bool
}

func (q GenericEntityQuery) Test(id uint16, em *EntityManager) bool {
	return q.TestFunc(id, em)
}

func GenericEntityQueryForTag(tag string) GenericEntityQuery {
	return GenericEntityQuery{
		Name: tag,
		TestFunc: func(id uint16, em *EntityManager) bool {
			return em.EntityHasTag(id, tag)
		}}
}

type EntityQueryWatcher struct {
	// the query this watcher will watch
	Query EntityQuery
	// A channel along which entity ID's will be sent, with the possibility
	// that those IDs are negative, with -(ID + 1) corresponding to ID
	// being deactivated
	Channel chan EntityToken
	Name    string
	// the ID of this watcher (used for memory management)
	ID uint16
}

// Construct a new entity query watcher (its channel will be created at the
// capacity of ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY constant)
func NewEntityQueryWatcher(
	q EntityQuery, name string, ID uint16) EntityQueryWatcher {

	return EntityQueryWatcher{
		q,
		make(chan EntityToken, ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY),
		name,
		ID}
}
