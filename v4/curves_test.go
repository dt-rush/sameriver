package sameriver

import (
	"math"
	"testing"
)

func TestCurves(t *testing.T) {
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
	}
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i/2], expect[i])
		}
	}
}
