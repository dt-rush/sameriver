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

func (c *curves) Lint(a float64, b float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		frac := (a - x) / (a - b)
		correction := c.Eq(1)(b)*c.Ge(b)(a) + c.Eq(0)(a)*c.Ge(b)(a)
		return math.Min(1, math.Max(0, frac)) - correction
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

func (c *curves) Eq(b float64) CurveFunc {
	return func(x float64) float64 {
		if x != b {
			return 0
		} else {
			return 1
		}
	}
}

// = 1 between [a,b)
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

func (c *curves) Peak(u float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		num := math.Abs(x - c.Gt(u)(x))
		denom := c.Gt(u)(x) - u
		return 1 - math.Sqrt(1-math.Pow(num/denom, 2))
	}
}

func (c *curves) Bump(w float64, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return c.Sigmoid((1-w)/2, s/5)(x) - c.Sigmoid((1+w)/2, s/5)(x)
	}
}

func (c *curves) Circ(x float64) float64 {
	x = c.Clamped(x)
	return math.Sqrt(-(2*x-1)*(2*x-1) + 1)
}

//
// PYRAMIDS
//

func (c *curves) Steps(n int, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		sum := 0.0
		for i := 1; i < n; i++ {
			sum += c.Sigmoid(0.5, s/5)(x + 0.5 - float64(i)/float64(n))
		}
		return sum / (float64(n) - 1)
	}
}

func (c *curves) StepsB(n int, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		sum := 0.0
		for i := 1; i < n; i++ {
			sum += c.Sigmoid(0.5, s/5)((1-1/float64(n))*x + 0.5 + 1/(2*float64(n)) - float64(i)/float64(n))
		}
		return sum / (float64(n) - 1)
	}
}

func (c *curves) Mayan(n int, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return c.StepsB(n, s)(2*x) - c.StepsB(n, s)(2*x-1)
	}
}

func (c *curves) SkewMayan(n int, u float64, s float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		d := (c.Lint(0, u)(x) + c.Lint(u, 1)(x)) / 2
		return c.Mayan(n, s)(d)
	}
}

func (c *curves) Pyramid(x float64) float64 {
	x = c.Clamped(x)
	return math.Max(0, 1-math.Abs(2*x-1))
}

func (c *curves) SkewPyramid(u float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		d := (c.Lint(0, u)(x) + c.Lint(u, 1)(x)) / 2
		return math.Max(0, 2*c.Tri(d/2)-1)
	}
}

//
// AUDIO
//

func (c *curves) Decay(k float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return math.Min(1, math.Pow(math.Abs(x-1), k)) * c.Le(1)(x)
	}
}

func (c *curves) Tri(x float64) float64 {
	x = c.Clamped(x)
	a := math.Max(0, (-math.Abs(4*x-1) + 1))
	b := math.Max(0, (-math.Abs(4*x-3) + 1))
	return (a - b + 1) / 2
}

func (c *curves) SqDuty(p float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return c.Lt(p)(math.Mod(x, 1))
	}
}

func (c *curves) Square(x float64) float64 {
	x = c.Clamped(x)
	return c.SqDuty(0.5)(x)
}

func (c *curves) Sin(x float64) float64 {
	x = c.Clamped(x)
	return (math.Sin(2*math.Pi*x) + 1) / 2
}

func (c *curves) LillyWave(x float64) float64 {
	x = c.Clamped(x)
	ratio := 34.0 / 28.0
	x1 := 2 * x
	a := math.Pow(c.Bell(0.5, 0.25)(x1), 2)
	x2 := ratio*x + 0.5 - ratio*0.75
	b := math.Pow(c.Bell(0.5, 0.25)(x2), 2)
	return (a - b + 1) / 2
}

func (c *curves) Comb(n int) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		param := float64(n) * math.Mod(x, 1/float64(n))
		return 0.5 * c.Circ(param)
	}
}

func (c *curves) Spring(freq float64, damp float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		base := func(x float64) float64 {
			return math.Exp(-damp*x) * math.Sin(2*math.Pi*freq*x)
		}
		norm := 1 / (base(math.Atan(2*math.Pi*freq/damp) / (2 * math.Pi * freq)))
		return (base(x)*norm + 1) / 2
	}
}

func (c *curves) SpringFlat(freq float64, damp float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		base := func(x float64) float64 {
			return math.Exp(-damp*x) * (math.Cos(2*math.Pi*freq*x) + 1) / 2
		}
		scale := (2*(freq-1) + 1) / (2 * freq)
		return base(scale * x)
	}
}

// BOUNCE
func (c *curves) Bounce(cor float64, tscale float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		B0 := func(x float64) float64 {
			return 1 - 0.5*9.81*tscale*x*tscale*x
		}
		rb0 := math.Sqrt(2*9.81*tscale*tscale) / (9.81 * tscale * tscale)
		root := func(k float64) float64 {
			return rb0 + 2*rb0*(cor*(1-math.Pow(cor, k))/(1-cor))
		}
		bounceNumber := func(x float64) float64 {
			return math.Ceil(math.Log(1-(1-cor)*(x/rb0-1)/(2*cor)) / math.Log(cor))
		}
		bn := bounceNumber(x)
		ex := math.Pow(cor, bn)
		arg := (x-c.Gt(0)(bn)*root(bn-1))/ex - c.Gt(0)(bn)*rb0
		return ex * B0(arg)
	}
}

func (c *curves) T(u float64) CurveFunc {
	return func(x float64) float64 {
		x = c.Clamped(x)
		return x
	}
}

func (c *curves) Clamped(x float64) float64 {
	return math.Min(1, math.Max(0, x))
}
