package sameriver

import (
	"math"
)

// visual aid:
// https://www.desmos.com/calculator/ylh21kqg4o

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

func (c *curves) Abs(x float64) float64 {
	return math.Abs(x)
}

func (c *curves) Linear(x float64) float64 {
	return (x + 1) / 2
}

func (c *curves) Exp(x float64) float64 {
	q := math.Exp(2) / (math.Exp(2) - 1)
	r := 1 / math.Exp(2)
	return q * (math.Exp(x-1) - r)
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
