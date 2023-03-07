package sameriver

import (
	"math"
)

// visual aid:
// https://www.desmos.com/calculator/ylh21kqg4o

type CurveFunc func(x float64) float64

// B(x, m_u, s_d) in desmos
// a simple bell curve
func Bell(mean, spread float64) CurveFunc {
	return func(x float64) float64 {
		return 1.0 / (((x-mean)*(x-mean))/(spread*spread) + 1)
	}
}

// S(x, s) in desmos
func Sigmoid(s float64) CurveFunc {
	return func(x float64) float64 {
		return math.Pow(1.0/(1+math.Exp(-x/s)), math.Abs(x-1))
	}
}

// F(x, m_u, s) in desmos
func Shelf(start, steepness float64) CurveFunc {
	return func(x float64) float64 {
		z := math.Pow((2 - start), 10*steepness)
		normalizeFactor := z / (z - 1)
		shelf := 1 - math.Pow(x-start+1, -10*steepness)
		return normalizeFactor * shelf
	}
}

// A_c(x, u, m) in desmos
// sqrt arc
func Sqrt(start, steepness float64) CurveFunc {
	return func(x float64) float64 {
		quot := math.Sqrt(x-start) / math.Sqrt(1-start)
		return math.Pow(quot, 1/steepness)
	}
}

// Q(x, k) in desmos
func QuadK(k float64) CurveFunc {
	return func(x float64) float64 {
		return math.Pow(x, 2*k)
	}
}

// Q(x, 1) in desmos
func Quad(x float64) float64 {
	return x * x
}

// A(x) in desmos
func Abs(x float64) float64 {
	return math.Abs(x)
}

// L(x) in desmos
func Linear(x float64) float64 {
	return (x + 1) / 2
}

// L_i(x, a, b) in desmos
func LinearInterval(a, b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, (x-a)/(b-a)))
	}
}

// P_k(x,p) in desmos represents the posPeak branch
func Peak(p float64) CurveFunc {
	// the positive half since we use sqrt
	posPeak := func(x, p float64) float64 {
		return math.Pow(1-math.Sqrt((x-p)/(1-p)), 2)
	}
	return func(x float64) float64 {
		if x >= p {
			return posPeak(x, p)
		} else {
			return posPeak(-x, -p)
		}
	}
}

// X(x) in desmos
func Exp(x float64) float64 {
	q := math.Exp(2) / (math.Exp(2) - 1)
	r := 1 / math.Exp(2)
	return q * (math.Exp(x-1) - r)
}

// G(x, b) in desmos
func Greater(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(x-b)))
	}
}

// E(x, b) in desmos
func Less(b float64) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, math.Ceil(b-x)))
	}
}

// I(x, a, b) in desmos
func Interval(a, b float64) CurveFunc {
	return func(x float64) float64 {
		return Greater(a)(x) + Less(b)(x) - 1
	}
}

// S_t(x, n, m) in desmos
func ContinuousStep(n int, smooth float64) CurveFunc {
	k := float64(n)
	m := smooth * 0.2
	normalizeFactor := (k + 1) / (k - 1)
	return func(x float64) float64 {
		y := 0.0
		for i := 1.0; i < k; i++ {
			step := Sigmoid(m)(
				10*(i/k+(x-x/k-1)/2),
			) / (k + 1)
			y += step
		}
		return normalizeFactor * y
	}
}

// C_m(x, n, s) in desmos
// a step pyramid centered at 1
func ContinuousMayan(n int, smooth float64) CurveFunc {
	return func(x float64) float64 {
		return ContinuousStep(n, smooth)(2*x+1) - ContinuousStep(n, smooth)(2*x-1)
	}
}

// S_p(x, n, m_u, s) in desmos
func SkewedMayan(n int, peak float64, smooth float64) CurveFunc {
	return func(x float64) float64 {
		// deform x from linear
		// such that it is -1 at x=-1, 0 at x=peak, 1 at x=1
		d := LinearInterval(-1, peak)(x) + LinearInterval(peak, 1)(x) - 1
		return ContinuousMayan(n, smooth)(d)
	}
}

// Q_u(x, n) in desmos
// quantize f into n bands in y (available y-values are only i*(1/n))
func QuantizeY(f CurveFunc, n int) CurveFunc {
	k := float64(n)
	return func(x float64) float64 {
		return math.Floor((f(x)+1/(2*(k-1)))*(k-1)) / (k - 1)
	}
}

// Q_x(x, n) in desmos (uses F_0 because functions aren't quite first class in desmos)
// quantize f into n bands in x (input to f is quantized i * (1/n))
func QuantizeX(f CurveFunc, n int) CurveFunc {
	k := float64(n)
	return func(x float64) float64 {
		// quantized x
		d := math.Floor(k*(x-1)/2)/(k/2) + 1 + 1/k
		return f(d)
	}
}

// D_c(x, k) in desmos
func Decay(k float64) CurveFunc {
	return func(x float64) float64 {
		return math.Pow(math.Abs(x-1)/2, k)
	}
}

// L_b (x, f_b) in desmos
func LinearBounce(f float64) CurveFunc {
	return func(x float64) float64 {
		return (1 - Linear(x)) * (math.Cos(f*math.Pi*(x+1)) + 1) / 2
	}
}

// E_b (x, f_b, k) in desmos
func ExponentialBounce(f float64, k float64) CurveFunc {
	return func(x float64) float64 {
		return Decay(k)(x) * (math.Cos(f*math.Pi*(x+1)) + 1) / 2
	}
}

func Clamped(f CurveFunc) CurveFunc {
	return func(x float64) float64 {
		return math.Min(1, math.Max(0, f(x)))
	}
}
