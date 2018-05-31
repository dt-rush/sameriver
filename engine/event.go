package engine

import (
	"fmt"
)

type Event struct {
	Type EventType
	Data interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("%s", EVENT_NAMES[e.Type])
}
