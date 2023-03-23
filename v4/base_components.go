package sameriver

const (
	POSITION ComponentID = iota
	VELOCITY
	ACCELERATION
	BOX
	MASS
	MAXVELOCITY
	BASESPRITE
	DESPAWNTIMER
	STEER
	MOVEMENTTARGET
	ITEM
	INVENTORY
	GENERICTAGS // NOTE: this should always be the last one, so clients can start
	// their consts at GENERICTAGS + 1 + iota
)
