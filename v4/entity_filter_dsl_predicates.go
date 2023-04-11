package sameriver

import (
	"strings"
)

// given the Identifier strings as args,
// the resolver strategy
// , we return a func(*Entity) bool, a "predicate", or "filter"
func EFDSLPredicatesBase(e *EntityFilterDSLEvaluator) EFDSLPredicateMap {

	return map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool){

		"Eq": e.Predicate(
			"any, any",
			func(k string, v int) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.HasComponent(STATE) && x.GetIntMap(STATE).ValCanBeSetTo(k, v)
				}
			},
		),

		"CanBe": e.Predicate(
			"string, int",
			func(k string, v int) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.HasComponent(STATE) && x.GetIntMap(STATE).ValCanBeSetTo(k, v)
				}
			},
		),

		"State": e.Predicate(
			"string, int",
			func(k string, v int) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.HasComponent(STATE) && x.GetIntMap(STATE).Get(k) == v
				}
			},
		),

		"HasComponent": e.Predicate(
			"string",
			func(componentStr string) func(*Entity) bool {
				return func(x *Entity) bool {
					// do a little odd access pattern since we only have
					// HasComponent for ComponentID (int) not strings.
					w := x.World
					ct := w.em.components
					componentID := ct.stringsRev[componentStr]
					return x.HasComponent(componentID)
				}
			},
		),

		"HasTag": e.Predicate(
			"string",
			func(tag string) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.HasTag(tag)
				}
			},
		),

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
				// distance from an entity
				entity := argsTyped[0].(*Entity)
				distance := argsTyped[1].(float64)
				return func(e *Entity) bool {
					pos := entity.GetVec2D(POSITION)
					box := entity.GetVec2D(BOX)
					return e.DistanceFromRect(*pos, *box) < distance
				}
			case 1:
				// distance from pos, box
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
		// TODO: withinpolygon, overlapspolygon

	}
}
