package sameriver

import (
	"testing"
)

func BenchmarkEFDSLSwitchIntFunc(b *testing.B) {
	// One-arg "int" func
	intFunc := func(arg int) func(*Entity) bool {
		return func(e *Entity) bool { return true }
	}

	// Prepare the arguments
	argsTyped := []any{42}

	for i := 0; i < b.N; i++ {
		_ = EFDSL.predicateSignatureAssertSwitch(intFunc, argsTyped)
	}
}

func BenchmarkEFDSLSwitchVec2DBoolSliceVec2DFunc(b *testing.B) {
	// Three-arg "*Vec2D, bool, []*Vec2D" func
	threeArgFunc := func(v *Vec2D, b bool, vs []*Vec2D) func(*Entity) bool {
		return func(e *Entity) bool { return true }
	}

	// Prepare the arguments
	argsTyped := []any{&Vec2D{1, 2}, true, []*Vec2D{{3, 4}, {5, 6}}}

	for i := 0; i < b.N; i++ {
		_ = EFDSL.predicateSignatureAssertSwitch(threeArgFunc, argsTyped)
	}
}
