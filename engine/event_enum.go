package engine

// TODO: generate

type EventType int

const N_EVENT_TYPES = 2

const (
	COLLISION_EVENT     = EventType(iota)
	SPAWN_REQUEST_EVENT = EventType(iota)
)

var EVENT_NAMES = map[EventType]string{
	COLLISION_EVENT:     "COLLISION_EVENT",
	SPAWN_REQUEST_EVENT: "SPAWN_REQUEST_EVENT",
}
