
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

type ComponentType int

const N_COMPONENT_TYPES = 4
const (
	BOX_COMPONENT      = iota
	SPRITE_COMPONENT   = iota
	TAGLIST_COMPONENT  = iota
	VELOCITY_COMPONENT = iota
)

var COMPONENT_NAMES = map[ComponentType]string{
	BOX_COMPONENT:      "BOX_COMPONENT",
	SPRITE_COMPONENT:   "SPRITE_COMPONENT",
	TAGLIST_COMPONENT:  "TAGLIST_COMPONENT",
	VELOCITY_COMPONENT: "VELOCITY_COMPONENT",
}