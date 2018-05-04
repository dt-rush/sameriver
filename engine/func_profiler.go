/**
  *
  * used to profile the execution time of functions
  *
  *
**/

package engine

import (
	"time"
)

const (
	FUNC_PROFILER_SIMPLE = iota
	FUNC_PROFILER_SIMPLE_MINMAX = iota
	FUNC_PROFILER_LOSSLESS = iota
)

// Accumulator data type to allow the easy computation of an average
type SimpleAccum struct {
	// total time consumed by all invocations
	TotalTime int
	// number of invocations
	Invocations int
}

// Accumulator data type to allow the easy computation of an average
// along with max, min values to give a slight sense of variance
type SimpleMinMaxAccum struct {
	// total time consumed by all invocations
	TotalTime int
	// number of invocations
	Invocations int
	// minimum and maximum values seen
	MinimumTime int
	MaximumTime int
}

// FuncProfiler stores data on function executions in one of 3 ways,
// depending on the mode given on Init
type FuncProfiler struct {
	// simple_data is a slice, indexed by the func ID, of 2-tuples of
	// (total ms, invocations) in order to find the average runtime
	simple_data []SimpleAccum
	// simple_minmax_data is a slice, indexed by the func ID, of 4-tuples of
	// (total ms, invocations, min, max) in order to find average
	// and extreme values
	simple_minmax_data []SimpleMinMaxAccum
	// lossless_data is a slice, indexed by the func ID, to a slice of
	// ints representing each runtime for each invocation.
	// NOTE:
	// This will use up a lot (relatively speaking) of memory if you aren't
	// careful. If the function is invoked on average every 50 ms, you'll
	// store 80 bytes per second, or 500 MB after about an hour and a half.
	// That might sound okay, but that's only one function. If you profile 4
	// such functions, you'll reach 500 MB after about 25 minutes. If you
	// profile 10 such functions, you'll reach 500 MB after 10 minutes. If you
	// profile a mere 6 functions invoked every 16 milliseconds, you'll reach
	// 500 MB in about 5 minutes. Really, this should only be used in a
	// specialized testing capacity to stress the engine and see what happens,
	// not for long-term tracking of statistics over the course of gameplay.
	lossless_data [][]int
	// start_times is a slice, indexed by func ID, of ints
	// used to time individual invocations. The user calls StartTimer(),
	// which puts the start time into this map, and in EndTimer(), the
	// current time is compared with the value stored here to determine the
	// runtime of the function
	start_times []int
	// mode tells which mode this FuncProfiler is working in
	Mode int
	// n_funcs keeps track of the number of functions which have
	// been registered for profiling
	n_funcs int
}

func (fp *FuncProfiler) Init(mode int) {
	// Initialize the data structures used to store profiling info
	fp.start_times = make ([]int, 1)
	switch mode {
		case FUNC_PROFILER_SIMPLE:
			fp.simple_data = make ([]SimpleAccum, 1)
		case FUNC_PROFILER_SIMPLE_MINMAX:
			fp.simple_minmax_data = make ([]SimpleMinMaxAccum, 1)
		case FUNC_PROFILER_LOSSLESS:
			fp.lossless_data = make ([][]int, 1)
	}
	// Set the mode
	fp.Mode = mode
	// We start at 0 funcs registered
	fp.n_funcs = 0
}

// register a function to be profiled. Allocates the mode-appropriate
// form of accumulator, returning the ID so the caller can interact with
// StartTimer and EndTimer
func (fp *FuncProfiler) Allocate() int {
	// generate the ID
	id := fp.n_funcs + 1
	// increment the number of funcs stored
	fp.n_funcs += 1
	// allocate storage for the profiling info based on the mode
	switch fp.Mode {
		case FUNC_PROFILER_SIMPLE:
			fp.simple_data = append (fp.simple_data,
				SimpleAccum{})
		case FUNC_PROFILER_SIMPLE_MINMAX:
			fp.simple_minmax_data = append (fp.simple_minmax_data,
				SimpleMinMaxAccum{})
		case FUNC_PROFILER_LOSSLESS:
			fp.lossless_data = append (fp.lossless_data,
				make ([]int, 1))
	}
	return id
}

// Start timing a function
func (fp *FuncProfiler) StartTimer (id int) {
	fp.start_times[id] = int (time.Now().UnixNano())
}

// Finish timing a function
func (fp *FuncProfiler) EndTimer (id int) {
	// get the elapsed milliseconds by comparing the current time to the
	// stored start time
	end_time := int (time.Now().UnixNano())
	milliseconds := (end_time - fp.start_times [id]) / 1e6
	// store the statistic according to the mode
	switch fp.Mode {
		case FUNC_PROFILER_SIMPLE:
			accum := &fp.simple_data [id]
			accum.TotalTime += milliseconds
			accum.Invocations += 1
		case FUNC_PROFILER_SIMPLE_MINMAX:
			accum := &fp.simple_minmax_data [id]
			accum.TotalTime += milliseconds
			accum.Invocations += 1
			if milliseconds < accum.MinimumTime {
				accum.MinimumTime = milliseconds
			}
			if milliseconds > accum.MaximumTime {
				accum.MaximumTime = milliseconds
			}
		case FUNC_PROFILER_LOSSLESS:
			fp.lossless_data[id] = append (fp.lossless_data[id],
				milliseconds)
	}
}
