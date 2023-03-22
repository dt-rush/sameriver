package sameriver

import (
	"testing"
)

func TestComponentBitArrayToString(t *testing.T) {
	w := testingWorld()
	b := w.em.components.BitArrayFromIDs([]ComponentID{POSITION, BOX, GENERICTAGS})
	s := w.em.components.BitArrayToString(b)
	// TODO: check s
	Logger.Printf(s)
}
