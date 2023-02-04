package sameriver

// state values are either bool or GOAPCtxState (bool resolved by get())
type GOAPState interface{}

var EmptyGOAPState = map[string]GOAPState{}

// a state val which can *set* modal values in the worldstate as the
// action chain runs forward
type GOAPCtxStateVal struct {
	name string
	val  bool
	get  func(ws GOAPWorldState) bool
	set  func(ws *GOAPWorldState)
}
