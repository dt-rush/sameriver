package utils

type TimeAccumulator struct {
	accum  float64
	period float64
}

func NewTimeAccumulator(period float64) TimeAccumulator {
	t := TimeAccumulator{}
	t.accum = 0
	t.period = period
	return t
}

func (t *TimeAccumulator) Tick(dt float64) bool {
	t.accum += dt
	had_tick := false
	for t.accum >= t.period {
		t.accum -= t.period
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.accum / t.period
}
