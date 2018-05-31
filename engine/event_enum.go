package engine

// TODO: generate

type EventType int

const N_EVENT_TYPES = 2

const (
	COLLISION_EVENT     = EventType(iota)
	SPAWN_REQUEST_EVENT = EventType(iota)
)

var EVENT_NAMES = map[EventType]string{
	COLLISION_EVENT:     "collision event",
	SPAWN_REQUEST_EVENT: "spawn request event",
}
