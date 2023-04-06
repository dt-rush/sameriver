package sameriver

var EntityFilterDSLSorts = map[string](func(args []string, resolver IdentifierResolver) func(a, b *Entity) int){
	"Closest": func(args []string, resolver IdentifierResolver) func(a, b *Entity) int {
		pole := resolver.Resolve(args[0]).(*Entity)
		return func(a, b *Entity) int {
			// float precision in an int package
			return int(10000*a.DistanceFrom(pole) - 10000*b.DistanceFrom(pole))
		}
	},
}
