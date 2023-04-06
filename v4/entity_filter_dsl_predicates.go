package sameriver

import "strconv"

var EntityFilterDSLPredicates = map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool){
	"CanBe": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		k, v := args[0], args[1]
		vi, err := strconv.Atoi(v)
		if err != nil {
			logWarning("CanBe got non-numeric argument in EntityFilter DSL: %s; will not behave properly. This must be fixed.", v)
		}
		return func(x *Entity) bool {
			return x.GetIntMap(STATE).ValCanBeSetTo(k, vi)
		}
	},
	"HasTag": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			return x.HasTag(args[0])
		}
	},
	"HasTags": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			return x.HasTags(args...)
		}
	},
	"Is": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		return func(x *Entity) bool {
			lookup := resolver.Resolve(args[0])
			if ent, ok := lookup.(*Entity); ok {
				return x == ent
			}
			return false
		}
	},
}
