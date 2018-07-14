package utils

type TimeAccumulator struct {
	Accum  float64
	Period float64
}

func CreateTimeAccumulator(period float64) TimeAccumulator {
	t := TimeAccumulator{}
	t.Accum = 0
	t.Period = period
	return t
}

func (t *TimeAccumulator) Tick(dt float64) bool {
	t.Accum += dt
	had_tick := false
	for t.Accum >= t.Period {
		t.Accum -= t.Period
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.Accum / t.Period
}
