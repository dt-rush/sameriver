package sameriver

func (e *EntityFilterDSLEvaluator) Predicate(signature string, f ...any) func(args []string, resolver IdentifierResolver) func(*Entity) bool {
	return func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, err := DSLAssertArgTypes(signature, args, resolver)
		if err != nil {
			logDSLError("%s", err)
			return nil
		}
		var result func(*Entity) bool
		// check if type signature is user-defined
		result = e.userPredicateSignatureAsserter(f[0], argsTyped)
		if result != nil {
			return result
		}
		// else, we handle a finite set of signatures
		return e.predicateSignatureAssertSwitch(f[0], argsTyped)
	}
}

func (e *EntityFilterDSLEvaluator) Sort(signature string, f ...any) func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) int {
	return func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) int {
		argsTyped, err := DSLAssertArgTypes(signature, args, resolver)
		if err != nil {
			logDSLError("%s", err)
			return nil
		}
		var result func(xs []*Entity) func(i, j int) int
		// check if type signature is user-defined
		result = e.userSortSignatureAsserter(f[0], argsTyped)
		if result != nil {
			return result
		}
		// else, we handle a finite set of signatures
		return e.sortSignatureAssertSwitch(f[0], argsTyped)

	}
}
