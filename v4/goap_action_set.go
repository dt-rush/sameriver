package sameriver

import (
	"fmt"
)

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
		name := action.DisplayName()
		if _, ok := as.set[name]; ok {
			panic(fmt.Sprintf("Action %s already added to actionset!", name))
		}
		as.set[name] = action
	}
}
