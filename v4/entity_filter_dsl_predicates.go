package sameriver

import (
	"strings"
)

var EntityFilterDSLPredicates = map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool){

	"CanBe": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("CanBe(string, int)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		k, v := argsTyped[0].(string), argsTyped[1].(int)
		return func(x *Entity) bool {
			return x.HasComponent(STATE) && x.GetIntMap(STATE).ValCanBeSetTo(k, v)
		}
	},
	"State": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("State(string, int)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		k, v := argsTyped[0].(string), argsTyped[1].(int)
		return func(x *Entity) bool {
			// TODO: instead of just args[1] = an int, what if its' a string
			// that can* be an int, but can also be like >3, or <=80
			return x.HasComponent(STATE) && x.GetIntMap(STATE).Get(k) == v
		}
	},
	"HasComponent": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("HasComponent(string)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		componentStr := strings.ToUpper(argsTyped[0].(string))
		return func(x *Entity) bool {
			w := x.World
			ct := w.em.components
			component := ct.stringsRev[componentStr]
			return x.HasComponent(component)
		}
	},
	"HasTag": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("HasTag(string)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		tag := argsTyped[0].(string)
		return func(x *Entity) bool {
			return x.HasTag(tag)
		}
	},
	"HasTags": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("HasTags([]string)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		tags := argsTyped[0].([]string)
		return func(x *Entity) bool {
			return x.HasTags(tags...)
		}
	},
	"Is": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("Is(IdentResolve<*Entity>)", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		entity := argsTyped[0].(*Entity)
		return func(x *Entity) bool {
			return x == entity
		}
	},
	"WithinDistance": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		signatures := []string{
			"WithinDistance(IdentResolve<*Entity>, float64)",
			"WithinDistance(IdentResolve<*Vec2D>, IdentResolve<*Vec2D>, float64)",
		}
		argsTyped, i, err := DSLAssertOverloadedArgTypes(signatures, args, resolver)
		if err != nil {
			logDSLError("WithinDistance: %s", err)
			logDSLError("Failed to match argument types for WithinDistance(%s)", strings.Join(args, ", "))
			return nil
		}
		switch i {
		case 0:
			entity := argsTyped[0].(*Entity)
			distance := argsTyped[1].(float64)
			return func(e *Entity) bool {
				pos := entity.GetVec2D(POSITION)
				box := entity.GetVec2D(BOX)
				return e.DistanceFromRect(*pos, *box) < distance
			}
		case 1:
			pos := argsTyped[0].(*Vec2D)
			box := argsTyped[1].(*Vec2D)
			distance := argsTyped[2].(float64)
			return func(e *Entity) bool {
				return e.DistanceFromRect(*pos, *box) < distance
			}
		default:
			logDSLError("Invalid arguments for WithinDistance, no matching signature found. Expected signatures: %v, got: WithinDistance(%s)", signatures, strings.Join(args, ", "))
			return nil
		}
	},
	"RectOverlap": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes("RectOverlap", args, resolver)
		if err != nil {
			logDSLError("%s", err)
		}
		pos := argsTyped[0].(*Vec2D)
		box := argsTyped[1].(*Vec2D)

		return func(e *Entity) bool {
			ePos := e.GetVec2D(POSITION)
			eBox := e.GetVec2D(BOX)

			return RectIntersectsRect(*pos, *box, *ePos, *eBox)
		}
	},
}
