package sameriver

func (e *EntityFilterDSLEvaluator) predicateSignatureAssertSwitch(f any, argsTyped []any) func(*Entity) bool {
	var result func(*Entity) bool
	switch fTyped := f.(type) {
	case func(int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int))
	case func(string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string))
	case func(string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int))
	case func(string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string))
	case func(int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int))
	default:
		panic("No case in either engine or user-supplied for the given func signature (use EFDSL.RegisterUserPredicateSignatureAsserter()")
	}

	return result
}

func (e *EntityFilterDSLEvaluator) sortSignatureAssertSwitch(f any, argsTyped []any) func(xs []*Entity) func(i, j int) int {
	var result func(xs []*Entity) func(i, j int) int
	switch fTyped := f.(type) {
	case func(int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int))
	case func(string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string))
	case func(string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int))
	case func(string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string))
	case func(int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int))
	default:
		panic("No case in either engine or user-supplied for the given func signature (use EFDSL.RegisterUserSortSignatureAsserter()")
	}

	return result
}
