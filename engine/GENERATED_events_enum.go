
//
//
// THIS FILE HAS BEEN GENERATED BY sameriver-generate
//
//
// DO NOT MODIFY BY HAND UNLESS YOU WANNA HAVE A GOOD TIME WHEN THE NEXT
// GENERATION DESTROYS WHAT YOU WROTE. UNLESS YOU KNOW HOW TO HAVE A GOOD TIME
//
//

package engine

type EventType int

const N_EVENT_TYPES = 2
const (
	COLLISION_EVENT = iota
	GENERIC_EVENT   = iota
)

var EVENT_NAMES = map[EventType]string{
	COLLISION_EVENT: "COLLISION_EVENT",
	GENERIC_EVENT:   "GENERIC_EVENT",
}