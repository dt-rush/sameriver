package sameriver

import (
	"strconv"
	"strings"
)

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
	// WithinDistance :: IdentResolve<*Entity.Vec2D> -> IdentResolve<*Entity.Vec2D> -> float64 -> *Entity -> bool
	"WithinDistance": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		// overloading hell yeah F U Rob Pike we smuggled it in
		signature := -1

		// ultimately we need two *Vec2D
		var posp, boxp *Vec2D

		// try WithinDistance(*Entity, float64)
		ent, entOk := resolver.Resolve(args[0]).(*Entity)
		d, atofErr := strconv.ParseFloat(args[1], 64)
		if entOk && atofErr == nil {
			signature = 0
		}
		posp = ent.GetVec2D(POSITION)
		boxp = ent.GetVec2D(BOX)

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
		boxp, boxOk := arg1.(*Vec2D)
		if !boxOk {
			var box Vec2D
			// (don't care if it's a pointer or the value itself)
			box, boxOk = arg0.(Vec2D)
			boxp = &box // re-pointer it lol
		}
		d, atofErr = strconv.ParseFloat(args[2], 64)
		if posOk && boxOk && atofErr == nil {
			signature = 1
		}

		if signature == -1 {
			logDSLError("WithinDistance() invocation was neither WithinDistance(*Vec2D,*Vec2D,float64) nor WithinDistance(*Entity,float64). Got WithinDistance(%args)")
			if atofErr != nil {
				logDSLError("%v", atofErr)
			}
		}

		return func(e *Entity) bool {
			return e.DistanceFromRect(*posp, *boxp) < d
		}
	},
}
