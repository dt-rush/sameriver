package sameriver

// EFDSLSorts key is the Function, like
// "Closest"
// given the Identifiers that were interspersed with commas as strings as args,
// and the resolver strategy,
// we return a func(xs []*Entity) func(i, j int) int, a "comparator/sorter"
// the return value of which, the func(i, j int) int with closure access to x,
// can be used to both compare / sort elements in xs.
// since the i, j func needs the closure reference to xs to actually sort it
// that's just the way sort.Slice() goes kid
func EFDSLSortsBase(e *EFDSLEvaluator) EFDSLSortMap {

	return map[string](func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool){

		"Closest": func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool {
			argsTyped, err := DSLAssertArgTypes("Closest(IdentResolve<*Entity>)", args, resolver)
			if err != nil {
				logDSLError("%s", err)
			}
			pole := argsTyped[0].(*Entity)
			return func(xs []*Entity) func(i, j int) bool {
				return func(i, j int) bool {
					return xs[i].DistanceFrom(pole) < xs[j].DistanceFrom(pole)
				}
			}
		},
	}
}
