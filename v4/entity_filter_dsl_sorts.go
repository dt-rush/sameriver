package sameriver

var EntityFilterDSLSorts = map[string](func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool){
	"Closest": func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool {
		pole := resolver.Resolve(args[0]).(*Entity)
		return func(xs []*Entity) func(i, j int) bool {
			return func(i, j int) bool {
				return xs[i].DistanceFrom(pole) < xs[j].DistanceFrom(pole)
			}
		}
	},
}
