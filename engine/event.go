package engine

import (
	"fmt"
)

type Event struct {
	Type        EventType
	Description string
	Data        interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("%s: %s", EVENT_NAMES[e.Type], e.Description)
}
