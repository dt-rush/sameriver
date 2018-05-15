/*
 *
 * Simple function profiler, recording total, count, min and max
 *
 */

package engine

import (
	"fmt"
)

// Accumulator data type to allow the easy computation of an average
// along with max, min values to give a slight sense of variance
type simpleMinMaxAccum struct {
	// total time consumed by all invocations
	totalTime float64
	// number of invocations
	invocations int
	// minimum and maximum values seen
	minimumTime float64
	maximumTime float64
}

// Create a new simpleMinMaxAccum
func newSimpleMinMaxAccum() simpleMinMaxAccum {
	return simpleMinMaxAccum{0, 0, 0, 0}
}

// Profiler struct using the simple-minmax accum
type simpleMinMaxProfiler struct {
	accum []simpleMinMaxAccum
	base  ProfilerBase
}

// Create a new instance of simpleMinMaxProfiler
func NewSimpleMinMaxProfiler() *simpleMinMaxProfiler {
	return &simpleMinMaxProfiler{
		accum: make([]simpleMinMaxAccum, 0),
		base:  NewProfilerBase(),
	}
}

// Register a function for profiling
func (p *simpleMinMaxProfiler) RegisterFunc(name string) uint16 {
	id := p.base.RegisterFunc(name)
	p.accum = append(p.accum, newSimpleMinMaxAccum())
	return id
}

// Clear the data associated with a function
func (p *simpleMinMaxProfiler) ClearData(id uint16) {
	p.accum[id] = newSimpleMinMaxAccum()
}

// Start timing a function for simpleMinMaxProfiler
func (p *simpleMinMaxProfiler) StartTimer(id uint16) {
	p.base.StartTimer(id)
}

// Stop timing a function for simpleMinMaxProfiler
func (p *simpleMinMaxProfiler) EndTimer(id uint16) {
	ms := p.base.EndTimer(id)
	accum := &p.accum[id]
	accum.totalTime += ms
	accum.invocations += 1
	if ms < accum.minimumTime {
		accum.minimumTime = ms
	}
	if ms > accum.maximumTime {
		accum.maximumTime = ms
	}
}

// Get the average runtime for a function
func (p *simpleMinMaxProfiler) GetAvg(id uint16) float64 {
	return (p.accum[id].totalTime /
		float64(p.accum[id].invocations))
}

// Get the name of a given function
func (p *simpleMinMaxProfiler) GetName(id uint16) string {
	return p.base.GetName(id)
}

// Set the name of a given function
func (p *simpleMinMaxProfiler) SetName(id uint16, name string) {
	p.base.SetName(id, name)
}

// Return a string displaying the stats for a given function
func (p *simpleMinMaxProfiler) GetSummaryString(id uint16) string {
	return fmt.Sprintf("Summary for %s: {Avg: %.3f}",
		p.base.GetName(id),
		p.GetAvg(id))
}
