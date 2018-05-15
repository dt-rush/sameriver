/*
 *
 * Simple function profiler, recording total and count
 *
 */

package engine

import (
	"fmt"
)

// Accumulator data type to allow the easy computation of an average
type simpleAccum struct {
	// total time consumed by all invocations
	totalTime float64
	// number of invocations
	invocations int
}

// Create a new simpleAccum
func newSimpleAccum() simpleAccum {
	return simpleAccum{0, 0}
}

// Profiler struct using the simple accum
type simpleProfiler struct {
	accum []simpleAccum
	base  ProfilerBase
}

// Create a new instance of simpleProfiler
func NewSimpleProfiler() *simpleProfiler {
	return &simpleProfiler{
		accum: make([]simpleAccum, 0),
		base:  NewProfilerBase(),
	}
}

// Register a function for profiling
func (p *simpleProfiler) RegisterFunc(name string) uint16 {
	id := p.base.RegisterFunc(name)
	p.accum = append(p.accum, newSimpleAccum())
	return id
}

// Clear the data associated with a function
func (p *simpleProfiler) ClearData(id uint16) {
	p.accum[id] = newSimpleAccum()
}

// Start timing a function
func (p *simpleProfiler) StartTimer(id uint16) {
	p.base.StartTimer(id)
}

// Stop timing a function
func (p *simpleProfiler) EndTimer(id uint16) {
	ms := p.base.EndTimer(id)
	accum := &p.accum[id]
	accum.totalTime += ms
	accum.invocations += 1
}

// Get the average runtime for a function
func (p *simpleProfiler) GetAvg(id uint16) float64 {
	return (p.accum[id].totalTime /
		float64(p.accum[id].invocations))
}

// Get the name of a given function
func (p *simpleProfiler) GetName(id uint16) string {
	return p.base.GetName(id)
}

// Set the name of a given function
func (p *simpleProfiler) SetName(id uint16, name string) {
	p.base.SetName(id, name)
}

// Return a string displaying the stats for a given function
func (p *simpleProfiler) GetSummaryString(id uint16) string {
	return fmt.Sprintf("Summary for %s: {Avg: %.3f}",
		p.base.GetName(id),
		p.GetAvg(id))
}
