/*
 *
 * Used to accumulate dt's (small increments of time), providing basic
 * utility functions based on this accumulation
 *
 */

package engine

// private struct
type TimeAccumulator struct {
	// how much time has accumulated so far
	accum int
	// the periodicity of the time accumulator (used by both Tick() and
	// Completion()
	period int
}

// Create a TimeAccumulator object with a given period
func NewTimeAccumulator(period int) TimeAccumulator {
	t := TimeAccumulator{}
	t.accum = 0
	t.period = period
	return t
}

// Add `dt` to the accumulated time and return true in that case that
// the accumulated time plus `dt` pushes us past the period.
// NOTE: will return a single true value even if `dt` is in fact
// greater than the period (an odd situation, but important to note). `dt`
// could be 100, and period could be 7, and we would still get a single true
// value, even though really 14 complete periods had elapsed.
func (t *TimeAccumulator) Tick(dt int) bool {
	t.accum += dt
	for t.accum >= t.period {
		t.accum %= t.period
		return true
	}
	return false
}

// Give the percent complete of the timer out of its period, given the
// current state of the accumulator
func (t *TimeAccumulator) Completion() float64 {
	return float64(t.accum) / float64(t.period)
}
