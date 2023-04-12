package sameriver

// given the Identifier strings as args,
// the resolver strategy
// , we return a func(*Entity) bool, a "predicate", or "filter"
func EFDSLPredicatesBase(e *EFDSLEvaluator) EFDSLPredicateMap {

	return map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool){

		"Eq": e.Predicate(
			"IdentResolve<any>, IdentResolve<any>",
			func(x any, y any) func(*Entity) bool {
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

		"HasTags": e.Predicate(
			"[]string",
			func(tags []string) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.HasTags(tags...)
				}
			},
		),

		"Is": e.Predicate(
			"IdentResolve<*Entity>",
			func(y *Entity) func(*Entity) bool {
				return func(x *Entity) bool {
					return x == y
				}
			},
		),

		"WithinDistance": e.Predicate(
			"IdentResolve<*Entity>, float64",
			func(y *Entity, d float64) func(*Entity) bool {
				return func(x *Entity) bool {
					pos := y.GetVec2D(POSITION)
					box := y.GetVec2D(BOX)
					return x.DistanceFromRect(*pos, *box) < d
				}
			},
			"IdentResolve<*Vec2D>, IdentResolve<*Vec2D>, float64",
			func(pos *Vec2D, box *Vec2D, d float64) func(*Entity) bool {
				return func(x *Entity) bool {
					return x.DistanceFromRect(*pos, *box) < d
				}
			},
		),

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
