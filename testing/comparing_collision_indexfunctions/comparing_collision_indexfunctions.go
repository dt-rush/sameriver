/* This was written to compare the time involved in computing unique indexes
 * given two uint16's via various methods.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// used to iterate a lot of times
const N = 1024 * 8

// a type representing a function which, given two uint16's, will return
// a uint32 index
type IndexFunc func(i uint16, j uint16) uint32

// iterate (N * (N-1)) / 2  times (that is, produce every pairing of numbers
// from 0 to N-1), and use the supplied IndexFunc to put an ExampleStruct
// into a map. return the milliseconds taken to do this
func computeABunchOfIndexes(index IndexFunc, name string) int64 {

	t0 := time.Now().UnixNano()

	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			// compute the index given the IndexFunc and store an object
			// in the map
			ix := index(uint16(i), uint16(j))
			// once in a while report a representative index to the console
			if rand.Intn(1e8) == 0 {
				fmt.Printf("representative index for %s(%d,%d): %d\n",
					name, i, j, ix)
			}
		}
	}

	t1 := time.Now().UnixNano()
	return (t1 - t0) / 1e6
}

// save the i part of the ix function in the i part of the loop
func computeABunchOfIndexesTriangleSpecial(index IndexFunc, name string) int64 {

	t0 := time.Now().UnixNano()

	for i := 0; i < N; i++ {
		ixoffset := -i*((i-1)/2-N) - (2*i + 1)
		for j := i + 1; j < N; j++ {
			ix := ixoffset + j
			// once in a while report a representative index to the console
			if rand.Intn(1e8) == 0 {
				fmt.Printf("representative index for %s(%d,%d): %d\n",
					name, i, j, ix)
			}
		}
	}

	t1 := time.Now().UnixNano()
	return (t1 - t0) / 1e6
}

// the number of times to run each call to computeABunchOfIndexes() in profile()
const N_TIMES = 8

// call computeABunchOfIndexes a certain number of times with the supplied IndexFunc
// and calculate the average time taken
func profile(
	indexFunc IndexFunc,
	indexFuncStr string,
	computeABunchOfIndexes func(IndexFunc, string) int64,
	name string) {

	durations := [N_TIMES]int64{}
	for i := 0; i < N_TIMES; i++ {
		durations[i] = computeABunchOfIndexes(indexFunc, name)
	}
	// calculate and report the average
	sum := int64(0)
	for i := 0; i < N_TIMES; i++ {
		sum += durations[i]
	}
	fmt.Println()
	fmt.Printf("avg time for `%s`: %d ms\n", indexFuncStr, sum/N_TIMES)
	fmt.Println()

}

func main() {

	rand.Seed(time.Now().UnixNano())

	// run the test with low indexes
	profile(func(i uint16, j uint16) uint32 {
		return uint32(i*N + j)
	},
		"i * N + j",
		computeABunchOfIndexes,
		"squareIndexes")

	// run the test with shifted indexes
	profile(func(i uint16, j uint16) uint32 {
		return uint32(i)<<16 | uint32(j)
	},
		"i << 16 | j",
		computeABunchOfIndexes,
		"shiftedIndexes")

	// run the test with triangle indexes
	triangle := func(n uint16) uint16 { return n * (n + 1) / 2 }
	profile(func(i uint16, j uint16) uint32 {
		return uint32((triangle(N) - triangle(N-i)) + j - 1 - (2*i + 1))
	},
		"(triangle(N) - triangle(N-i)) + j - 1 - (2*i + 1)",
		computeABunchOfIndexes,
		"triangleIndexes")

	// run the test with triangle indexes, simplified by algebra
	profile(func(i uint16, j uint16) uint32 {
		return uint32(-i*((i-1)/2-N) + j - (2*i + 1))
	},
		"-i*((i-1)/2-N) + j - (2*i + 1)",
		computeABunchOfIndexesTriangleSpecial,
		"triangleSimplifiedIndexes")
}
