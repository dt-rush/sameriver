package engine

type testSystem struct {
	x float64
}

func newTestSystem() *testSystem {
	return &testSystem{}
}
func (s *testSystem) LinkWorld(w *World) {}
func (s *testSystem) Update(dt_ms float64) {
	s.x += dt_ms
}

type testDependentSystem struct {
	ts *testSystem `sameriver-system-dependency:"-"`
}

func newTestDependentSystem() *testDependentSystem {
	return &testDependentSystem{}
}
func (s *testDependentSystem) LinkWorld(w *World)   {}
func (s *testDependentSystem) Update(dt_ms float64) {}
