package engine

import (
	"fmt"
)

type EventType int

type Event struct {
	Type        EventType
	Description string
	Data        interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("%d:%s", e.Type, e.Description)
}

// NOTE: the below constants must be kept in line with the structs
// to allow receivers of game events to type assert their events properly
// in order to unwrap the data inside, and to allow the game event manager
// to work properly
// TODO: assert this is correct during build
const N_GAME_EVENT_TYPES = 2
const (
	EVENT_TYPE_COLLISION     = iota
	EVENT_TYPE_SPAWN_REQUEST = iota
)

type CollisionEventData struct {
	EntityA uint16
	EntityB uint16
}

type SpawnRequestData struct {
	EntityType int
	Position   [2]int16
	Active     bool
}
