package engine

type GameEvent struct {
	Class       int
	Description string
	Data        interface{}
}

func (e GameEvent) String() string {
	return fmt.Sprintf("%d:%s", e.Type, e.Description)
}

// NOTE: the below constants must be kept in line with the structs
// to allow receivers of game events to type assert their events properly
// in order to unrwap the data inside
const (
	COLLISION_EVENT = iota
)

type CollisionEvent struct {
	EntityA uint16
	EntityB uint16
}
