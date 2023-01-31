package utils

type TimeAccumulator struct {
	accum_ms  float64
	period_ms float64
}

func NewTimeAccumulator(period_ms float64) TimeAccumulator {
	t := TimeAccumulator{}
	t.accum_ms = 0
	t.period_ms = period_ms
	return t
}

func (t *TimeAccumulator) Tick(dt_ms float64) bool {
	Logger.Printf("t.accum_ms: %f, t.period: %f, dt_ms: %f", t.accum_ms, t.period_ms, dt_ms)
	t.accum_ms += dt_ms
	had_tick := false
	for t.accum_ms >= t.period_ms {
		t.accum_ms -= t.period_ms
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.accum_ms / t.period_ms
}
