package sameriver

type Blackboard struct {
	Name   string
	State  map[string]any
	Events *EventBus
}

func NewBlackboard(name string) *Blackboard {
	return &Blackboard{
		Name:   name,
		State:  make(map[string]any),
		Events: NewEventBus("blackboard-" + name),
	}
}
