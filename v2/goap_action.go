package sameriver

type GOAPAction struct {
	name string
	// either an int or a func() int
	cost interface{}
	pres map[string]GOAPState
	effs map[string]GOAPState
}

func (a *GOAPAction) presFulfilled(ws GOAPWorldState) bool {
	state := NewGOAPWorldState(nil)
	for name, val := range a.pres {
		if ctxVal, ok := val.(GOAPCtxStateVal); ok {
			if !ctxVal.get(ws) {
				return false
			}
		}
		state.Vals[name] = val
	}
	return ws.fulfills(state)
}
