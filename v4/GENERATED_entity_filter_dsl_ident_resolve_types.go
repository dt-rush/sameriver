package sameriver

var IdentResolveTypeAssertMap = map[string]DSLArgTypeAssertionFunc{
	"*Entity": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Entity](resolver.Resolve(arg), "*Entity")
	},
	"string": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[string](resolver.Resolve(arg), "string")
	},
	"*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Vec2D](resolver.Resolve(arg), "*Vec2D")
	},
	"[]*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[[]*Vec2D](resolver.Resolve(arg), "[]*Vec2D")
	},
	"*EventPredicate": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*EventPredicate](resolver.Resolve(arg), "*EventPredicate")
	},
	// Add more types here...
}
