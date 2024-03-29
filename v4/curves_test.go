package sameriver

import (
	"math"
	"testing"
)

func TestCurvesClimb(t *testing.T) {
	expect := []float64{
		Curves.Sigmoid(0.5, 1)(-1), 0,
		Curves.Sigmoid(0.5, 1)(0), 0,
		Curves.Sigmoid(0.5, 1)(0.5), 0.5,
		Curves.Sigmoid(0.5, 1)(1), 1,
		Curves.Sigmoid(0.5, 1)(2), 1,

		Curves.Shelf(0.5, 1)(-1), 0,
		Curves.Shelf(0.5, 1)(0), 0,
		Curves.Shelf(0.5, 1)(0.5), 0,
		Curves.Shelf(0.5, 1)(1), 1,
		Curves.Shelf(0.5, 1)(2), 1,

		Curves.Shelf(1, 1)(-1), 0,
		Curves.Shelf(1, 1)(0), 0,
		Curves.Shelf(1, 1)(0.5), 0,
		Curves.Shelf(1, 1)(1), 1,
		Curves.Shelf(1, 1)(2), 1,

		Curves.Quad(0.5, 1)(-1), 0,
		Curves.Quad(0.5, 1)(0), 0,
		Curves.Quad(0.5, 1)(0.5), 0,
		Curves.Quad(0.5, 1)(1), 1,
		Curves.Quad(0.5, 1)(2), 1,

		Curves.Cubi(0.5, 1)(-1), 0,
		Curves.Cubi(0.5, 1)(0), 0,
		Curves.Cubi(0.5, 1)(0.5), 0,
		Curves.Cubi(0.5, 1)(1), 1,
		Curves.Cubi(0.5, 1)(2), 1,

		Curves.Lin(-1), 0,
		Curves.Lin(0), 0,
		Curves.Lin(0.5), 0.5,
		Curves.Lin(1), 1,
		Curves.Lin(2), 1,

		Curves.Lint(0.3, 0.7)(-1), 0,
		Curves.Lint(0.3, 0.7)(0), 0,
		Curves.Lint(0.3, 0.7)(0.25), 0,
		Curves.Lint(0.3, 0.7)(0.4), 0.25,
		Curves.Lint(0.3, 0.7)(0.5), 0.5,
		Curves.Lint(0.3, 0.7)(0.6), 0.75,
		Curves.Lint(0.3, 0.7)(0.75), 1,
		Curves.Lint(0.3, 0.7)(1), 1,
		Curves.Lint(0.3, 0.7)(2), 1,

		Curves.Exp(5)(-1), 0,
		Curves.Exp(5)(0), 0,
		Curves.Exp(5)(1), 1,
		Curves.Exp(5)(2), 1,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesPeaks(t *testing.T) {
	expect := []float64{
		Curves.Bell(0.5, 1)(-1), 0.058,
		Curves.Bell(0.5, 1)(0), 0.058,
		Curves.Bell(0.5, 1)(0.5), 1,
		Curves.Bell(0.5, 1)(1), 0.058,
		Curves.Bell(0.5, 1)(2), 0.058,

		Curves.BellPinned(0.5)(-1), 0,
		Curves.BellPinned(0.5)(0), 0,
		Curves.BellPinned(0.5)(0.5), 1,
		Curves.BellPinned(0.5)(1), 0,
		Curves.BellPinned(0.5)(2), 0,

		Curves.Plateau(4)(-1), 0,
		Curves.Plateau(4)(0), 0,
		Curves.Plateau(4)(0.5), 1,
		Curves.Plateau(4)(1), 0,
		Curves.Plateau(4)(2), 0,

		Curves.Peak(0.5)(-1), 0,
		Curves.Peak(0.5)(0), 0,
		Curves.Peak(0.5)(0.5), 1,
		Curves.Peak(0.5)(1), 0,
		Curves.Peak(0.5)(2), 0,

		Curves.Bump(0.5, 1)(-1), 0,
		Curves.Bump(0.5, 1)(0), 0,
		Curves.Bump(0.5, 1)(0.5), 1,
		Curves.Bump(0.5, 1)(1), 0,
		Curves.Bump(0.5, 1)(2), 0,

		Curves.Circ(-1), 0,
		Curves.Circ(0), 0,
		Curves.Circ(0.5), 1,
		Curves.Circ(1), 0,
		Curves.Circ(2), 0,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesPyramids(t *testing.T) {
	expect := []float64{
		Curves.Steps(5)(-1), 0,
		Curves.Steps(5)(0), 0,
		Curves.Steps(5)(0.1), 0,
		Curves.Steps(5)(0.3), 0.25,
		Curves.Steps(5)(0.5), 0.5,
		Curves.Steps(5)(0.7), 0.75,
		Curves.Steps(5)(0.9), 1,
		Curves.Steps(5)(1), 1,
		Curves.Steps(5)(2), 1,

		Curves.StepsB(5)(-1), 0,
		Curves.StepsB(5)(0), 0,
		Curves.StepsB(5)(0.1), 0,
		Curves.StepsB(5)(0.14), 0.25,
		Curves.StepsB(5)(0.4), 0.5,
		Curves.StepsB(5)(0.644), 0.75,
		Curves.StepsB(5)(0.7), 0.75,
		Curves.StepsB(5)(0.85), 1,
		Curves.StepsB(5)(0.9), 1,
		Curves.StepsB(5)(1), 1,
		Curves.StepsB(5)(2), 1,

		Curves.StepsCont(5, 0.2)(-1), 0,
		Curves.StepsCont(5, 0.2)(0), 0,
		Curves.StepsCont(5, 0.2)(0.1), 0,
		Curves.StepsCont(5, 0.2)(0.3), 0.25,
		Curves.StepsCont(5, 0.2)(0.5), 0.5,
		Curves.StepsCont(5, 0.2)(0.7), 0.75,
		Curves.StepsCont(5, 0.2)(0.9), 1,
		Curves.StepsCont(5, 0.2)(1), 1,
		Curves.StepsCont(5, 0.2)(2), 1,

		Curves.StepsBCont(5, 0.1)(-1), 0,
		Curves.StepsBCont(5, 0.1)(0), 0,
		Curves.StepsBCont(5, 0.1)(0.1), 0,
		Curves.StepsBCont(5, 0.1)(0.14), 0.25,
		Curves.StepsBCont(5, 0.1)(0.4), 0.5,
		Curves.StepsBCont(5, 0.1)(0.6), 0.5,
		Curves.StepsBCont(5, 0.1)(0.7), 0.75,
		Curves.StepsBCont(5, 0.1)(0.85), 0.75,
		Curves.StepsBCont(5, 0.1)(0.9), 1,
		Curves.StepsBCont(5, 0.1)(1), 1,
		Curves.StepsBCont(5, 0.1)(2), 1,

		Curves.Mayan(5, 0.2)(-1), 0,
		Curves.Mayan(5, 0.2)(0), 0,
		Curves.Mayan(5, 0.2)(0.1), 0.25,
		Curves.Mayan(5, 0.2)(0.25), 0.5,
		Curves.Mayan(5, 0.2)(0.4), 0.75,
		Curves.Mayan(5, 0.2)(0.47), 1,
		Curves.Mayan(5, 0.2)(0.53), 1,
		Curves.Mayan(5, 0.2)(0.6), 0.75,
		Curves.Mayan(5, 0.2)(0.75), 0.5,
		Curves.Mayan(5, 0.2)(0.85), 0.25,
		Curves.Mayan(5, 0.2)(1), 0,
		Curves.Mayan(5, 0.2)(2), 0,

		Curves.SkewMayan(5, 0.2, 0.2)(-1), 0,
		Curves.SkewMayan(5, 0.2, 0.2)(0), 0,
		Curves.SkewMayan(5, 0.2, 0.2)(0.048), 0.25,
		Curves.SkewMayan(5, 0.2, 0.2)(0.1), 0.5,
		Curves.SkewMayan(5, 0.2, 0.2)(0.15), 0.75,
		Curves.SkewMayan(5, 0.2, 0.2)(0.193), 1,
		Curves.SkewMayan(5, 0.2, 0.2)(0.275), 1,
		Curves.SkewMayan(5, 0.2, 0.2)(0.3777), 0.75,
		Curves.SkewMayan(5, 0.2, 0.2)(0.5413), 0.5,
		Curves.SkewMayan(5, 0.2, 0.2)(0.757), 0.25,
		Curves.SkewMayan(5, 0.2, 0.2)(0.94), 0,
		Curves.SkewMayan(5, 0.2, 0.2)(1), 0,
		Curves.SkewMayan(5, 0.2, 0.2)(2), 0,

		Curves.Pyramid(-1), 0,
		Curves.Pyramid(0), 0,
		Curves.Pyramid(0.5), 1,
		Curves.Pyramid(1), 0,
		Curves.Pyramid(2), 0,

		Curves.SkewPyramid(0.2)(-1), 0,
		Curves.SkewPyramid(0.2)(0), 0,
		Curves.SkewPyramid(0.2)(0.2), 1,
		Curves.SkewPyramid(0.2)(1), 0,
		Curves.SkewPyramid(0.2)(2), 0,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesAudio(t *testing.T) {
	expect := []float64{
		Curves.Decay(2)(-1), 1,
		Curves.Decay(2)(0), 1,
		Curves.Decay(2)(0.5), 0.25,
		Curves.Decay(2)(0.75), 0.0625,
		Curves.Decay(2)(1), 0,
		Curves.Decay(2)(2), 0,

		Curves.Tri(-1), 0.5,
		Curves.Tri(0), 0.5,
		Curves.Tri(0.25), 1,
		Curves.Tri(0.5), 0.5,
		Curves.Tri(0.75), 0,
		Curves.Tri(1), 0.5,
		Curves.Tri(2), 0.5,

		Curves.SqDuty(0.25)(-1), 1,
		Curves.SqDuty(0.25)(0), 1,
		Curves.SqDuty(0.25)(0.24), 1,
		Curves.SqDuty(0.25)(0.26), 0,
		Curves.SqDuty(0.25)(0.5), 0,
		Curves.SqDuty(0.25)(0.75), 0,

		Curves.Square(-1), 1,
		Curves.Square(0), 1,
		Curves.Square(0.49), 1,
		Curves.Square(0.51), 0,
		Curves.Square(0.75), 0,

		Curves.Sin(-1), 0.5,
		Curves.Sin(0), 0.5,
		Curves.Sin(0.25), 1,
		Curves.Sin(0.5), 0.5,
		Curves.Sin(0.75), 0,
		Curves.Sin(1), 0.5,
		Curves.Sin(2), 0.5,

		Curves.LillyWave(-1), 0.5,
		Curves.LillyWave(0), 0.5,
		Curves.LillyWave(0.25), 1,
		Curves.LillyWave(0.5), 0.5,
		Curves.LillyWave(0.75), 0,
		Curves.LillyWave(1), 0.5,
		Curves.LillyWave(2), 0.5,

		Curves.Comb(4)(-1), 0,
		Curves.Comb(4)(0), 0,
		Curves.Comb(4)(0.125), 0.5,
		Curves.Comb(4)(0.25), 0,
		Curves.Comb(4)(0.375), 0.5,
		Curves.Comb(4)(0.5), 0,
		Curves.Comb(4)(0.625), 0.5,
		Curves.Comb(4)(0.75), 0,
		Curves.Comb(4)(0.875), 0.5,
		Curves.Comb(4)(1), 0,
		Curves.Comb(4)(2), 0,

		Curves.Spring(19, 13.7)(-1), 0.5,
		Curves.Spring(19, 13.7)(0), 0.5,
		Curves.Spring(19, 13.7)(0.0648), 0.7431,
		Curves.Spring(19, 13.7)(0.2227), 0.5279,
		Curves.Spring(19, 13.7)(0.2227), 0.5279,
		Curves.Spring(19, 13.7)(1), 0.5,
		Curves.Spring(19, 13.7)(2), 0.5,

		Curves.SpringFlat(19, 13.7)(-1), 1,
		Curves.SpringFlat(19, 13.7)(0), 1,
		Curves.SpringFlat(19, 13.7)(0.0521), 0.4927,
		Curves.SpringFlat(19, 13.7)(0.3224), 0.0134,
		Curves.SpringFlat(19, 13.7)(1), 0,
		Curves.SpringFlat(19, 13.7)(2), 0,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesBounce(t *testing.T) {
	c := 0.377902785645
	expect := []float64{
		Curves.Bounce(c, 1.0)(-1), 1,
		Curves.Bounce(c, 1.0)(0), 1,
		Curves.Bounce(c, 1.0)(0.622), 0.3779,
		Curves.Bounce(c, 1.0)(1), 0,
		Curves.Bounce(c, 1.0)(2), 0,

		Curves.NBounce(c, 3)(-1), 1,
		Curves.NBounce(c, 3)(0), 1,
		Curves.NBounce(c, 3)(0.4899), 0,
		Curves.NBounce(c, 3)(0.675), 0.3779,
		Curves.NBounce(c, 3)(0.8601), 0,
		Curves.NBounce(c, 3)(0.93), 0.1428,
		Curves.NBounce(c, 3)(1), 0,
		Curves.NBounce(c, 3)(2), 0,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesQuantize(t *testing.T) {
	expect := []float64{
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(-1), 0,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(0), 0,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(0.24), 0,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(0.3), 0.1429,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(0.5), 0.5714,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(0.76), 1,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(1), 1,
		Curves.QuantY(8, Curves.Sigmoid(0.5, 2))(2), 1,

		Curves.QuantX(16, Curves.Circ)(-1), 0,
		Curves.QuantX(16, Curves.Circ)(0), 0,
		Curves.QuantX(16, Curves.Circ)(0.07), 0.4989,
		Curves.QuantX(16, Curves.Circ)(0.245), 0.8,
		Curves.QuantX(16, Curves.Circ)(0.5), 0.9978,
		Curves.QuantX(16, Curves.Circ)(0.78), 0.8,
		Curves.QuantX(16, Curves.Circ)(0.9), 0.4989,
		Curves.QuantX(16, Curves.Circ)(1), 0,
		Curves.QuantX(16, Curves.Circ)(2), 0,
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}
