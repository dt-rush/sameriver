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

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type ComponentSet struct {
	Box            *Vec2D
	Logic          *LogicUnit
	Mass           *float64
	MaxVelocity    *float64
	MovementTarget *Vec2D
	Position       *Vec2D
	Sprite         *Sprite
	Steer          *Vec2D
	TagList        *TagList
	Velocity       *Vec2D
}

func (cs *ComponentSet) ToBitArray() bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENT_TYPES))
	if cs.Box != nil {
		b.SetBit(uint64(BOX_COMPONENT))
	}
	if cs.Logic != nil {
		b.SetBit(uint64(LOGIC_COMPONENT))
	}
	if cs.Mass != nil {
		b.SetBit(uint64(MASS_COMPONENT))
	}
	if cs.MaxVelocity != nil {
		b.SetBit(uint64(MAXVELOCITY_COMPONENT))
	}
	if cs.MovementTarget != nil {
		b.SetBit(uint64(MOVEMENTTARGET_COMPONENT))
	}
	if cs.Position != nil {
		b.SetBit(uint64(POSITION_COMPONENT))
	}
	if cs.Sprite != nil {
		b.SetBit(uint64(SPRITE_COMPONENT))
	}
	if cs.Steer != nil {
		b.SetBit(uint64(STEER_COMPONENT))
	}
	if cs.TagList != nil {
		b.SetBit(uint64(TAGLIST_COMPONENT))
	}
	if cs.Velocity != nil {
		b.SetBit(uint64(VELOCITY_COMPONENT))
	}
	return b
}

func (em *EntityManager) ApplyComponentSet(cs ComponentSet) func(*EntityToken) {
	b := cs.ToBitArray()
	return func(entity *EntityToken) {
		// make sure the component bitarray matches the applied component set
		entity.ComponentBitArray = entity.ComponentBitArray.Or(b)
		if cs.Box != nil {
			em.Components.Box[entity.ID] = *cs.Box
		}
		if cs.Logic != nil {
			em.Components.Logic[entity.ID] = *cs.Logic
		}
		if cs.Mass != nil {
			em.Components.Mass[entity.ID] = *cs.Mass
		}
		if cs.MaxVelocity != nil {
			em.Components.MaxVelocity[entity.ID] = *cs.MaxVelocity
		}
		if cs.MovementTarget != nil {
			em.Components.MovementTarget[entity.ID] = *cs.MovementTarget
		}
		if cs.Position != nil {
			em.Components.Position[entity.ID] = *cs.Position
		}
		if cs.Sprite != nil {
			em.Components.Sprite[entity.ID] = *cs.Sprite
		}
		if cs.Steer != nil {
			em.Components.Steer[entity.ID] = *cs.Steer
		}
		if cs.TagList != nil {
			em.Components.TagList[entity.ID] = *cs.TagList
		}
		if cs.Velocity != nil {
			em.Components.Velocity[entity.ID] = *cs.Velocity
		}
	}
}
