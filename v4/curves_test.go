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
		Curves.Steps(5, 0.2)(-1), 0,
		Curves.Steps(5, 0.2)(0), 0,
		Curves.Steps(5, 0.2)(0.1), 0,
		Curves.Steps(5, 0.2)(0.3), 0.25,
		Curves.Steps(5, 0.2)(0.5), 0.5,
		Curves.Steps(5, 0.2)(0.7), 0.75,
		Curves.Steps(5, 0.2)(0.9), 1,
		Curves.Steps(5, 0.2)(1), 1,
		Curves.Steps(5, 0.2)(2), 1,

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
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}
