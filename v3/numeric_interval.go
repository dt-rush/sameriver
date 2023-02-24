package sameriver

import (
	"math"
)

type NumericInterval struct {
	A, B float64
}

// given an open interval [a, b], return the least amount needed
// to modify x such that it will be in bounds (0 if in bounds, + if below,
// - if above)
//
// NOTE: we have to use float64 for +,- Inf, which is float64 type... annoying
func (i *NumericInterval) Diff(x float64) float64 {
	if x >= i.A && x <= i.B {
		return 0
	} else if x < i.A {
		return i.A - x
	} else {
		return i.B - x
	}
}

func MakeNumericInterval(op string, val int) *NumericInterval {
	switch op {
	case "<":
		return &NumericInterval{math.Inf(-1), float64(val - 1)}
	case "<=":
		return &NumericInterval{math.Inf(-1), float64(val)}
	case "=":
		return &NumericInterval{float64(val), float64(val)}
	case ">=":
		return &NumericInterval{float64(val), math.Inf(+1)}
	case ">":
		return &NumericInterval{float64(val + 1), math.Inf(+1)}
		/*
			case ">;<":
				// TODO
			case ">=;<":
				// TODO
			case ">;<=":
				// TODO
			case ">=;<=":
				// TODO
		*/
	default:
		panic("Got undefined op in GOAPGoalFunc() [valid: >=,>,=,<,<=]")
	}
}
