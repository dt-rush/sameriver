package sameriver

type GOAPActionSet struct {
	set map[string]*GOAPAction
}

func NewGOAPActionSet() *GOAPActionSet {
	return &GOAPActionSet{
		set: make(map[string]*GOAPAction),
	}
}

func (as *GOAPActionSet) Add(actions ...*GOAPAction) {
	for _, action := range actions {
		as.set[action.name] = action
	}
}
