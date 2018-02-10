package engine

type TimeAccumulator struct {
    Accum int
    Period int
}

func CreateTimeAccumulator (period int) TimeAccumulator {
    t := TimeAccumulator{}
    t.Accum = 0
    t.Period = period
    return t
}

func (t *TimeAccumulator) Tick (dt int) bool {
    t.Accum += dt
    had_tick := false
    for t.Accum >= t.Period {
        t.Accum %= t.Period
        had_tick = true
    }
    return had_tick
}

func (t *TimeAccumulator) Completion () float64 {
    return float64 (t.Accum) / float64 (t.Period)
}
