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
	// Logger.Printf("    thoseThatHelpFulfill %s", ws.Vals)
	for _, action := range as.set {
		effState := NewGOAPWorldState(nil)
		effState = effState.applyAction(action)
		// Logger.Printf("    effState of %s:", action.name)
		// Logger.Printf("    %s", effState.Vals)
		if effState.isSubset(ws) {
			helpers.Add(action)
		}
	}
	return helpers
}
