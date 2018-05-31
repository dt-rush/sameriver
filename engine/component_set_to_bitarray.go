package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

// TODO: generate

func (s *ComponentSet) ToBitArray() bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENT_TYPES))
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
	if s.Mind != nil {
		b.SetBit(uint64(MIND_COMPONENT))
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
