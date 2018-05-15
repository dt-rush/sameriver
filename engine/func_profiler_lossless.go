/*
 *
 * Lossless profiler, recording every single ms runtime of functions
 *
 */

package engine

import (
	"fmt"
)

// Accumulator data type to allow lossless storage of invocation times,
// as a basis for detailed statistical computations
type losslessAccum struct {
	// times taken by each invocation as a simple data series
	// NOTE:
	// This will use up a lot (relatively speaking) of memory if you aren't
	// careful. If the function is invoked on average every 100 ms, you'll
	// store 80 bytes per second, or 500 MB after about an hour and a half.
	// That might sound okay, but that's only one function. If you profile 2
	// such functions, you'll reach 500 MB after about 25 minutes. If you
	// profile 5 such functions, you'll reach 500 MB after 10 minutes. If you
	// profile a mere 3 functions invoked every 16 milliseconds, you'll reach
	// 500 MB in about 5 minutes. Really, this should only be used in a
	// specialized testing capacity to stress the engine and see what happens,
	// not for long-term tracking of statistics over the course of gameplay.
	// TODO (add as feature instead of leaving up to the user):
	// Keeping the above in mind, there may be ways to save memory by adding
	// logic to periodically "reduce" the raw data to its statistical measures
	// and storing those, maybe even writing the raw data thus-reduced to disk
	// if requested.
	times []float64
}

// Create a new losslessAccum
func newLosslessAccum() losslessAccum {
	return losslessAccum{make([]float64, 0)}
}

// Profiler struct using the lossless accum
type losslessProfiler struct {
	accum []losslessAccum
	base  ProfilerBase
}

// Create a new instance of losslessProfiler
func NewLosslessProfiler() *losslessProfiler {
	return &losslessProfiler{
		accum: make([]losslessAccum, 0),
		base:  NewProfilerBase(),
	}
}

// Register a function for profiling with losslessProfiler
func (p *losslessProfiler) RegisterFunc(name string) uint16 {
	id := p.base.RegisterFunc(name)
	p.accum = append(p.accum, newLosslessAccum())
	return id
}

// Clear the data associated with a function for losslessProfiler
func (p *losslessProfiler) ClearData(id uint16) {
	p.accum[id] = newLosslessAccum()
}

// Start timing a function for losslessProfiler
func (p *losslessProfiler) StartTimer(id uint16) {
	p.base.StartTimer(id)
}

// Stop timing a function for losslessProfiler
func (p *losslessProfiler) EndTimer(id uint16) {
	ms := p.base.EndTimer(id)
	times := &p.accum[id].times
	*times = append(*times, ms)
}

// Get the average runtime for a function
func (p *losslessProfiler) GetAvg(id uint16) float64 {
	sum := 0.0
	count := len(p.accum[id].times)
	for i := 0; i < count; i++ {
		sum += p.accum[id].times[i]
	}
	return (sum / float64(count))
}

// Get the name of a given function
func (p *losslessProfiler) GetName(id uint16) string {
	return p.base.GetName(id)
}

// Set the name of a given function
func (p *losslessProfiler) SetName(id uint16, name string) {
	p.base.SetName(id, name)
}

// Return a string displaying the stats for a given function
func (p *losslessProfiler) GetSummaryString(id uint16) string {
	// TODO: add more statistics
	return fmt.Sprintf("Summary for %s: {Avg: %.3f}",
		p.base.GetName(id),
		p.GetAvg(id))
}
