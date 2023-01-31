package utils

type TimeAccumulator struct {
	accum     float64
	period_ms float64
}

func NewTimeAccumulator(period_ms float64) TimeAccumulator {
	t := TimeAccumulator{}
	t.accum = 0
	t.period_ms = period_ms
	return t
}

func (t *TimeAccumulator) Tick(dt float64) bool {
	Logger.Println("In Tick()")
	t.accum += dt
	had_tick := false
	for t.accum >= t.period_ms {
		Logger.Println("t.accum >= t.period_ms")
		t.accum -= t.period_ms
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.accum / t.period_ms
}
