package engine

// TODO: implement this interface in a struct which allows generic predication
// on entities (their component values and anything else)
type EntityQuery interface {
	Test(id uint16, entity_manager *EntityManager) bool
}

type EntityQueryWatcher struct {
	Query   EntityQuery
	Channel chan (int16)
	Name    string
}
