package sameriver

type GOAPActionSet struct {
	set map[string]GOAPAction
}

func NewGOAPActionSet() *GOAPActionSet {
	return &GOAPActionSet{
		set: make(map[string]GOAPAction),
	}
}

func (as *GOAPActionSet) Add(actions ...GOAPAction) {
	for _, action := range actions {
		as.set[action.name] = action
	}
}

func (as *GOAPActionSet) thoseThatHelpFulfill(ws GOAPWorldState) *GOAPActionSet {
	helpers := NewGOAPActionSet()
	for _, action := range as.set {
		effState := NewGOAPWorldState(nil)
		effState = effState.applyAction(action)
		if effState.isSubset(ws) {
			helpers.Add(action)
		}
	}
	return helpers
}
