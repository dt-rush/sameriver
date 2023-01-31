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

func (t *TimeAccumulator) Tick(dt_ms float64) bool {
	Logger.Println("In Tick()")
	t.accum += dt_ms
	had_tick := false
	Logger.Printf("t.accum: %f", t.accum)
	Logger.Printf("t.period_ms: %f", t.period_ms)
	for t.accum >= t.period_ms {
		Logger.Println("t.accum >= t.period_ms")
		t.accum -= t.period_ms
		Logger.Printf("after sub: t.accum: %f", t.accum)
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.accum / t.period_ms
}
