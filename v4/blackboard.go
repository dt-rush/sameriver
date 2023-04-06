package sameriver

type Blackboard struct {
	Name   string
	state  map[string]any
	Events *EventBus
}

func NewBlackboard(name string) *Blackboard {
	return &Blackboard{
		Name:   name,
		state:  make(map[string]any),
		Events: NewEventBus("blackboard-" + name),
	}
}

func (b *Blackboard) Has(k string) bool {
	_, ok := b.state[k]
	return ok
}

func (b *Blackboard) Get(k string) any {
	return b.state[k]
}

func (b *Blackboard) Set(k string, v any) {
	b.state[k] = v
}
