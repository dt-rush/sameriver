package sameriver

import (
	"math"
)

// visual aid:
// https://www.desmos.com/calculator/kw8fhm0qox

type curves struct{}

var Curves = curves{}

type CurveFunc func(x float64) float64

func (c *curves) Bell(spread float64) CurveFunc {
	return func(x float64) float64 {
		return 1.0 / ((x*x)/(spread*spread) + 1)
	}
}

func (c *curves) Sigmoid(x float64) float64 {
	return 1.0 / (1 + math.Exp(-x/0.2))
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
