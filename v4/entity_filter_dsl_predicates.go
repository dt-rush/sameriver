package sameriver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrDSLEntityAccessFailure = errors.New("value specified is not *Entity")
var ErrDSLExpectedTypeFailure = errors.New("identifier doesn't resolve to type wanted")

var EntityFilterDSLPredicates = map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool){

	// CanBe :: string -> int -> *Entity -> bool
	"CanBe": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		k, v := args[0], args[1]
		vi, err := strconv.Atoi(v)
		if err != nil {
			logDSLError("CanBe got non-numeric argument in EntityFilter DSL: %s; will not behave properly.", v)
		}
		return func(x *Entity) bool {
			return x.HasComponent(STATE) && x.GetIntMap(STATE).ValCanBeSetTo(k, vi)
		}
	},
	// State :: string -> int -> *Entity -> bool
	"State": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		k, v := args[0], args[1]
		vi, err := strconv.Atoi(v)
		if err != nil {
			logDSLError("State got non-numeric argument in EntityFilter DSL: %s; will not behave properly.", v)
		}
		return func(x *Entity) bool {
			return x.HasComponent(STATE) && x.GetIntMap(STATE).Get(k) == vi
		}
	},
	// HasComponent :: string -> *Entity -> bool
	"HasComponent": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		componentStr := strings.ToUpper(args[0])
		return func(x *Entity) bool {
			w := x.World
			ct := w.em.components
			component := ct.stringsRev[componentStr]
			return x.HasComponent(component)
		}
	},
	// HasTag :: string -> *Entity -> bool
	"HasTag": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			return x.HasTag(args[0])
		}
	},
	// HasTags :: []string -> *Entity -> bool
	// do we have the tags
	"HasTags": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			return x.HasTags(args...)
		}
	},
	// Is :: IdentResolve<*Entity> -> *Entity -> bool
	// are we a certain identity by pointer
	"Is": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			lookup := resolver.Resolve(args[0])
			if ent, ok := lookup.(*Entity); ok {
				return x == ent
			}
			return false
		}
	},
	// WithinDistance :: IdentResolve<*Entity> -> float64 -> *Entity -> bool
	// OR
	// WithinDistance :: IdentResolve<*Vec2D> -> IdentResolve<*Vec2D> -> float64 -> *Entity -> bool
	"WithinDistance": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		// overloading hell yeah F U Rob Pike we smuggled it in

		// ultimately we need two *Vec2D and a float,
		// so a signature is some way of transforming
		// our args into posp, boxp, d
		signatures := []func() (posp, boxp *Vec2D, d float64, err error){
			func() (posp, boxp *Vec2D, d float64, err error) {
				// try WithinDistance(*Entity, float64)
				ent, entOk := resolver.Resolve(args[0]).(*Entity)
				d, err = strconv.ParseFloat(args[1], 64)
				if err != nil {
					return nil, nil, -1, err
				}
				if entOk {
					posp = ent.GetVec2D(POSITION)
					boxp = ent.GetVec2D(BOX)
					return posp, boxp, d, nil
				} else {
					return nil, nil, -1, fmt.Errorf("%w for %s", ErrDSLEntityAccessFailure, args[0])
				}
			},
			func() (posp, boxp *Vec2D, d float64, err error) {
				// try WithinDistance(Vec2D, Vec2D, float64)
				arg0 := resolver.Resolve(args[0])
				arg1 := resolver.Resolve(args[1])
				posp, posOk := arg0.(*Vec2D)
				if !posOk {
					// (don't care if it's a pointer or the value itself)
					var pos Vec2D
					pos, posOk = arg0.(Vec2D)
					posp = &pos // re-pointer it lol
				}
				if !posOk {
					return nil, nil, -1, fmt.Errorf("%w for \"%s\" should be *Vec2D", ErrDSLExpectedTypeFailure, args[0])
				}
				boxp, boxOk := arg1.(*Vec2D)
				if !boxOk {
					var box Vec2D
					// (don't care if it's a pointer or the value itself)
					box, boxOk = arg0.(Vec2D)
					boxp = &box // re-pointer it lol
				}
				if !boxOk {
					return nil, nil, -1, fmt.Errorf("%w for \"%s\" should be *Vec2D", ErrDSLExpectedTypeFailure, args[1])
				}
				d, err = strconv.ParseFloat(args[2], 64)
				if err != nil {
					return nil, nil, -1, err
				}
				return posp, boxp, d, nil
			},
		}

		for _, sig := range signatures {
			posp, boxp, d, err := sig()
			if err == nil {
				return func(e *Entity) bool {
					return e.DistanceFromRect(*posp, *boxp) < d
				}
			}
		}
		logDSLError("WithinDistance() invocation was neither WithinDistance(*Vec2D,*Vec2D,float64) nor WithinDistance(*Entity,float64). Got WithinDistance(%s)", strings.Join(args, ","))
		return nil
	},
}
