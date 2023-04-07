package sameriver

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityPredicate func(*Entity) bool

// just for fun, let's define a useless function. let's be luxuriant about
// modern storage space. we can waste the bytes.
func NullEntityPredicate(e *Entity) bool {
	return false
}
func AllEntityPredicate(e *Entity) bool {
	return true
}

type EntityFilter struct {
	Name      string
	Predicate func(e *Entity) bool
}

func NewEntityFilter(
	name string, f func(e *Entity) bool) EntityFilter {
	return EntityFilter{Name: name, Predicate: f}
}

func (q EntityFilter) Test(e *Entity) bool {
	return q.Predicate(e)
}

func EntityFilterFromTag(tag string) EntityFilter {
	return EntityFilter{
		Name: tag,
		Predicate: func(e *Entity) bool {
			return e.GetTagList(GENERICTAGS).Has(tag)
		}}
}

func EntityFilterFromComponentBitArray(
	name string, q bitarray.BitArray) EntityFilter {
	return EntityFilter{
		Name: name,
		Predicate: func(e *Entity) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := e.ComponentBitArray
			return q.Equals(q.And(b))
		}}
}

func EntityFilterFromCanBe(canBe map[string]int) EntityFilter {
	return EntityFilter{
		Name: "canbe",
		Predicate: func(e *Entity) bool {
			for k, v := range canBe {
				if !e.GetIntMap(STATE).ValCanBeSetTo(k, v) {
					return false
				}
			}
			return true
		},
	}
}

// bit of a meta filter:
// matches the closest entity to to that fulfills the given filter
func EntityFilterFromClosest(to *Entity, filter EntityFilter) EntityFilter {
	return EntityFilter{
		Name: "closest",
		Predicate: func(e *Entity) bool {
			return e == to.World.ClosestEntityFilter(
				*to.GetVec2D(POSITION),
				*to.GetVec2D(BOX),
				filter.Predicate)
		},
	}
}
