package sameriver

// note: we break DRY here because it's much clearer to read what is going on. I tried to
// generalize the core of these and there were type nightmares that started to make it way worse
// in terms of readability

/*
Predicate() and Sort() are factory functions that simplify the definition of
EFDSL predicates and sorts.

They take an arbitrary number of interleaved string signatures and
corresponding functions as arguments. The factory function handles type
assertion, identifier resolution, and uses generated switch statements to pass
the type-asserted arguments to the corresponding function based on the provided
signature. This makes defining EFDSL predicates more concise and readable.
*/

func (e *EFDSLEvaluator) Predicate(args ...any) func(args []string, resolver IdentifierResolver) func(*Entity) bool {
	if len(args)%2 != 0 {
		panic("Mismatched signatures and functions in Predicate()")
	}

	signatures := make([]string, len(args)/2)
	funcs := make([]any, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		signatures[i/2] = args[i].(string)
		funcs[i/2] = args[i+1]
	}

	return func(args []string, resolver IdentifierResolver) func(*Entity) bool {
		argsTyped, i, err := DSLAssertOverloadedArgTypes(signatures, args, resolver)
		if err != nil {
			logDSLError("%s", err)
			return nil
		}

		// check if type signature is user-defined
		result := e.userPredicateSignatureAsserter(funcs[i], argsTyped)
		if result != nil {
			return result
		}
		// else, we handle a finite set of signatures
		return e.predicateSignatureAssertSwitch(funcs[i], argsTyped)
	}
}

func (e *EFDSLEvaluator) Sort(args ...any) func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) int {
	if len(args)%2 != 0 {
		panic("Mismatched signatures and functions in Sort()")
	}

	signatures := make([]string, len(args)/2)
	funcs := make([]any, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		signatures[i/2] = args[i].(string)
		funcs[i/2] = args[i+1]
	}

	return func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) int {
		argsTyped, i, err := DSLAssertOverloadedArgTypes(signatures, args, resolver)
		if err != nil {
			logDSLError("%s", err)
			return nil
		}

		// check if type signature is user-defined
		result := e.userSortSignatureAsserter(funcs[i], argsTyped)
		if result != nil {
			return result
		}
		// else, we handle a finite set of signatures
		return e.sortSignatureAssertSwitch(funcs[i], argsTyped)
	}
}
