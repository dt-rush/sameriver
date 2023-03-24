package sameriver

import (
	"math"
)

// visual aid:
// https://www.desmos.com/calculator/kw8fhm0qox

type curves struct{}

var Curves = curves{}

type CurveFunc func(x float64) float64

//
// CLIMB
//

func (c *curves) Sigmoid(u float64, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		denom := 1.0 + math.Exp(-(x-u)/(0.05*s))
		return math.Pow(1.0/denom, math.Abs((x-u)-1))
	}
}

func (c *curves) Shelf(u float64, s float64) CurveFunc {
	n0 := func(x0, s float64) float64 {
		num := math.Pow(2-x0, 10*s)
		denom := math.Pow(2-x0, 10*s) - 1
		return num / denom
	}
	return func(x float64) float64 {
		x = c.Clamped(x)
		if x == 1 {
			return 1
		}
		return math.Max(0, n0(u, s)*(1-math.Pow((x-u)+1, -10*s)))
	}
}

func (c *curves) Arc(u float64, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		if x == 1 {
			return 1
		}
		if x >= u {
			num := math.Sqrt(x - u)
			denom := math.Sqrt(1 - u)
			return math.Pow(num/denom, 1/s)
		} else {
			return 0
		}
	}
}

func (c *curves) Quad(u float64, k float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		if x == 1 {
			return 1
		}
		if x <= u {
			return 0
		} else {
			return math.Pow(x-u, 2*k) / math.Pow(1-u, 2*k)
		}
	}
}

func (c *curves) Cubi(u float64, k float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		if x == 1 {
			return 1
		}
		if x <= u {
			return 0
		} else {
			return (math.Pow((2/(1-u))*(x-u)-1, 3*k) + 1) / 2
		}
	}
}

func (c *curves) Lin(x float64) float64 {
	x = c.Clamped(x)
	return x
}

func (c *curves) Lint(u float64, w float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return math.Min(1, math.Max(0, (x-u+w)/(2*w)))
	}
}

func (c *curves) Exp(k float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		num := math.Exp(k*x) - 1
		denom := math.Exp(k) - 1
		return num / denom
	}
}

//
// INTERVALS
//

func (c *curves) Gt(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(x-b)))
	}
}

func (c *curves) Ge(b float64) CurveFunc {
	return func(x float64) float64 {
		return 1 - c.Lt(b)(x)
	}
}

func (c *curves) Lt(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(b-x)))
	}
}

func (c *curves) Le(b float64) CurveFunc {
	return func(x float64) float64 {
		return 1 - c.Gt(b)(x)
	}
}

// = 1 between a, b inclusive
// = 0 elsewhere
func (c *curves) Span(a, b float64) CurveFunc {
	return func(x float64) float64 {
		return c.Ge(a)(x) + c.Le(b)(x) - 1
	}
}

// PEAKS
func (c *curves) Bell(u float64, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return 1.0 / ((x-u)*(x-u)/(s*s/64) + 1)
	}
}

func (c *curves) BellPinned(u float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return math.Max(0, -16*math.Pow(x-(u-0.5), 3)*4*math.Pow((x-(u-0.5))-1, 3))
	}
}

func (c *curves) Plateau(k float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return math.Max(0, 1-math.Pow(2*(x-0.5), 2*(k+1)))
	}
}

func (c *curves) Clamped(x float64) float64 {
	return math.Min(1, math.Max(0, x))
}
