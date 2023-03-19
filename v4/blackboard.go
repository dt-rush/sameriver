package sameriver

type Blackboard struct {
	Name   string
	State  map[string]interface{}
	Events *EventBus
}

func NewBlackboard(name string) *Blackboard {
	return &Blackboard{
		Name:   name,
		State:  make(map[string]interface{}),
		Events: NewEventBus("blackboard-" + name),
	}
}
