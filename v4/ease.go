package sameriver

func Ease(curve CurveFunc, duration_ms float64, x *float64) (logic func(dt_ms float64)) {
	tick := NewTimeAccumulator(duration_ms)
	return func(dt_ms float64) {
		*x = curve(tick.CompletionAfterDT(dt_ms))
		tick.Tick(dt_ms)
	}
}
