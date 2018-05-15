package engine


// TODO: use to expand entity manager watchers to allow for watching of
// component values
type Query interface {
	Test(id uint16, entity_manager *EntityManager) bool
}

type QueryWatcher struct {
	Query Query
	Channel chan (int16)
	ID uint16
}
