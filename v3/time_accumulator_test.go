package sameriver

import (
	"fmt"
	"math"
	"testing"
)

func TestTimeAccumulatorTicks(t *testing.T) {
	testCases := [][3]float64{
		[3]float64{10, 3, 4},
		[3]float64{100, 30, 4},
		[3]float64{100, 1, 100},
	}
	for _, tcase := range testCases {
		ta := NewTimeAccumulator(tcase[0])
		i := 0.0
		hadTick := false
		for !hadTick {
			hadTick = ta.Tick(tcase[1])
			i++
		}
		if i != tcase[2] {
			t.Fatal("did not tick properly")
		}
	}
}

func TestTimeAccumulatorCompletion(t *testing.T) {
	testCases := [][4]float64{
		[4]float64{10, 0, 0, 0},
		[4]float64{10, 1, 10, 0.0},
		[4]float64{10, 1, 2, 0.2},
		[4]float64{10, 1, 5, 0.5},
		[4]float64{20, 2, 5, 0.5},
		[4]float64{30, 1, 15, 0.5},
		[4]float64{42, 21, 1, 0.5},
		[4]float64{10, 1, 9, 0.9},
		[4]float64{100, 1, 97, 0.97},
		[4]float64{1000, 1, 975, 0.975},
	}
	for _, tcase := range testCases {
		ta := NewTimeAccumulator(tcase[0])
		for i := 0.0; i < tcase[2]; i++ {
			ta.Tick(tcase[1])
		}
		completion := ta.Completion()
		if math.Abs(completion-tcase[3]) > 0.01 {
			t.Fatal(fmt.Sprintf("did not calculate completion properly for "+
				"%v", tcase))
		}
	}
}
