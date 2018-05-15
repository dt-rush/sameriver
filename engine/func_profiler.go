/**
  *
  * common definitions for the func profiler implementing structs
  * used to profile the execution time of functions
  *
  *
**/

package engine

import (
	"time"
)

const (
	FUNC_PROFILER_SIMPLE        = iota
	FUNC_PROFILER_SIMPLE_MINMAX = iota
	FUNC_PROFILER_LOSSLESS      = iota
)

// struct containing data members common to all profilers
type ProfilerBase struct {
	// start_times is a slice, indexed by func ID, of ints
	// used to time individual invocations. The user calls StartTimer(),
	// which puts the start time into this map, and in EndTimer(), the
	// current time is compared with the value stored here to determine the
	// runtime of the function
	start_times []int
	// names is a slice, indexed by func ID, of strings
	// which are the names of the accumulators
	names []string
	// n_funcs keeps track of the number of functions which have
	// been registered for profiling
	n_funcs int
}

// Create and return a new ProfilerBase struct
// with its data structures allocated
func NewProfilerBase() ProfilerBase {
	return ProfilerBase{
		start_times: make([]int, 1),
		names:       make([]string, 1),
		n_funcs:     0,
	}
}

// Do the bookkeeping to register a function in the ProfilerBase
func (b *ProfilerBase) RegisterFunc(name string) uint16 {
	// generate the ID
	id := uint16(b.n_funcs)
	// increment the number of funcs stored
	b.n_funcs += 1
	// store the name
	b.names = append(b.names, name)
	// create an entry in start times
	b.start_times = append(b.start_times, 0)
	// return the allocated ID
	return id
}

// Get the name of a function from a ProfilerBase object by ID
func (b *ProfilerBase) GetName(id uint16) string {
	return b.names[id]
}

// Set the name of a function from a ProfilerBase object by ID
func (b *ProfilerBase) SetName(id uint16, name string) {
	b.names[id] = name
}

// start timing a function with ProfilerBase
func (b *ProfilerBase) StartTimer(id uint16) {
	b.start_times[id] = int(time.Now().UnixNano())
}

// start timing a function with ProfilerBase, returning milliseconds elapsed
func (b *ProfilerBase) EndTimer(id uint16) float64 {
	end_time := int(time.Now().UnixNano())
	milliseconds := float64(end_time-b.start_times[id]) / float64(1e6)
	return milliseconds
}

// Generic methods of a FuncProfiler
type FuncProfiler interface {
	RegisterFunc(name string) uint16
	ClearData(id uint16)
	StartTimer(id uint16)
	EndTimer(id uint16)
	GetName(id uint16) string
	SetName(id uint16, name string)
	GetAvg(id uint16) float64
	GetSummaryString(id uint16) string
}

// Create and return a FuncProfiler given the requested mode
func NewFuncProfiler(mode int) FuncProfiler {
	var profiler FuncProfiler
	switch mode {
	case FUNC_PROFILER_SIMPLE:
		profiler = NewSimpleProfiler()
	case FUNC_PROFILER_SIMPLE_MINMAX:
		profiler = NewSimpleMinMaxProfiler()
	case FUNC_PROFILER_LOSSLESS:
		profiler = NewLosslessProfiler()
	}
	return profiler
}
