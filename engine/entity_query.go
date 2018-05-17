package engine

// TODO: use to expand entity manager watchers to allow for watching of
// component values
type EntityQuery interface {
	Test(id uint16, entity_manager *EntityManager) bool
}

type EntityQueryWatcher struct {
	Query   EntityQuery
	Channel chan (int16)
	Name    string
}
