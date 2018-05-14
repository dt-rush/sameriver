/*
 *
 * A set of component values which may be provided or not
 *
 */

package component

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

const (
	ACTIVE_COMPONENT   = iota
	COLOR_COMPONENT    = iota
	HITBOX_COMPONENT   = iota
	LOGIC_COMPONENT    = iota
	POSITION_COMPONENT = iota
	SPRITE_COMPONENT   = iota
	VELOCITY_COMPONENT = iota
)

const N_COMPONENTS = 7

type ComponentSet struct {
	Active   *bool
	Color    *sdl.Color
	Hitbox   *[2]uint16
	Logic    *LogicUnit
	Position *[2]uint16
	Sprite   *Sprite
	Velocity *[2]uint16
}

func (s *ComponentSet) ToBitarray() bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENTS))
	if s.Active != nil {
		b.SetBit(ACTIVE_COMPONENT)
	}
	if s.Color != nil {
		b.SetBit(COLOR_COMPONENT)
	}
	if s.Hitbox != nil {
		b.SetBit(HITBOX_COMPONENT)
	}
	if s.Logic != nil {
		b.SetBit(LOGIC_COMPONENT)
	}
	if s.Position != nil {
		b.SetBit(POSITION_COMPONENT)
	}
	if s.Sprite != nil {
		b.SetBit(SPRITE_COMPONENT)
	}
	if s.Velocity != nil {
		b.SetBit(VELOCITY_COMPONENT)
	}
	return b
}

func MakeComponentQuery(components []int) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENTS))
	for _, COMPONENT := range components {
		b.SetBit(COMPONENT)
	}
	return b
}
