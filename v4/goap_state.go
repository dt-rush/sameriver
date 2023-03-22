package sameriver

type IntOrFunc any

// a state val which can *set* modal values in the worldstate as the
// action chain runs forward, and which can *read* modal values when
// an action's pre is being checked
type GOAPModalVal struct {
	// the var name this value will be used with
	name string
	// a func that is used when this state val appears in a pre
	// eg. "atTree,=": atTreeModal
	// where ctxAtTree.valAsPre will look in the modal state for the entity
	// position and return 0 if not at the tree, 1 if at the tree
	//
	// also used to *re-evaluate* a state var's value during applyAction
	// (eg. if we had a modal var involving atTree, it will set atTree: 0
	// when our modal position has moved away from the tree)
	check func(ws *GOAPWorldState) int
	// a func that is used when this state val appears in an eff, modifying
	// modal state
	effModalSet func(ws *GOAPWorldState, op string, x int)
}
