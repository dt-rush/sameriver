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
		denominator := 1.0 + math.Exp(-(x-u)/(0.05*s))
		return math.Pow(1.0/denominator, math.Abs((x-u)-1))
	}
}

func (c *curves) Shelf(u float64, s float64) CurveFunc {
	n0 := func(x0, s float64) float64 {
		numerator := math.Pow(2-x0, 10*s)
		denominator := math.Pow(2-x0, 10*s) - 1
		return numerator / denominator
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
			numerator := math.Sqrt(x - u)
			denominator := math.Sqrt(1 - u)
			return math.Pow(numerator/denominator, 1/s)
		} else {
			return 0
		}
	}
}

func (c *curves) Bell(s float64) CurveFunc {
	return func(x float64) float64 {
		return 1.0 / ((x*x)/(s*s) + 1)
	}
}

func (c *curves) Quadratic(x float64) float64 {
	return x * x
}

func (c *curves) Linear(x float64) float64 {
	return (x + 1) / 2
}

func (c *curves) Greater(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(x-b)))
	}
}

func (c *curves) Less(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(b-x)))
	}
}

func (c *curves) Interval(a, b float64) CurveFunc {
	return func(x float64) float64 {
		return c.Greater(a)(x) + c.Less(b)(x) - 1
	}
}

func (c *curves) Clamped(x float64) float64 {
	return math.Min(1, math.Max(0, x))
}
