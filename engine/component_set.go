/*
 *
 * A set of component values which may be provided or not
 *
 */

package engine

import (
	"bytes"
	"fmt"

	"github.com/golang-collections/go-datastructures/bitarray"
	"github.com/veandco/go-sdl2/sdl"
)

//
// NOTE: if you wish to add a new component, you'll need to modify this file,
// and components_table.go
//
const (
	ACTIVE_COMPONENT   = iota
	COLOR_COMPONENT    = iota
	HEALTH_COMPONENT   = iota
	HITBOX_COMPONENT   = iota
	LOGIC_COMPONENT    = iota
	POSITION_COMPONENT = iota
	SPRITE_COMPONENT   = iota
	TAGLIST_COMPONENT  = iota
	VELOCITY_COMPONENT = iota
)

const N_COMPONENTS = 7

type ComponentSet struct {
	Active   *bool
	Color    *sdl.Color
	Health   *uint8
	HitBox   *[2]uint16
	Logic    *LogicUnit
	Position *[2]int16
	Sprite   *Sprite
	TagList  *TagList
	Velocity *[2]float32
}

func (s *ComponentSet) ToBitArray() bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENTS))
	if s.Active != nil {
		b.SetBit(uint64(ACTIVE_COMPONENT))
	}
	if s.Color != nil {
		b.SetBit(uint64(COLOR_COMPONENT))
	}
	if s.Health != nil {
		b.SetBit(uint64(HEALTH_COMPONENT))
	}
	if s.HitBox != nil {
		b.SetBit(uint64(HITBOX_COMPONENT))
	}
	if s.Logic != nil {
		b.SetBit(uint64(LOGIC_COMPONENT))
	}
	if s.Position != nil {
		b.SetBit(uint64(POSITION_COMPONENT))
	}
	if s.Sprite != nil {
		b.SetBit(uint64(SPRITE_COMPONENT))
	}
	if s.TagList != nil {
		b.SetBit(uint64(TAGLIST_COMPONENT))
	}
	if s.Velocity != nil {
		b.SetBit(uint64(VELOCITY_COMPONENT))
	}
	return b
}

func MakeComponentBitArray(components []int) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENTS))
	for _, COMPONENT := range components {
		b.SetBit(uint64(COMPONENT))
	}
	return b
}

func ComponentBitArrayToString(b bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < N_COMPONENTS; i++ {
		bit, _ := b.GetBit(i)
		var val int
		if bit {
			val = 1
		} else {
			val = 0
		}
		buf.WriteString(fmt.Sprintf("%v", val))
	}
	buf.WriteString("]")
	return buf.String()
}
