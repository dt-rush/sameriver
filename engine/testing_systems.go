package engine

// a basic system with a data member
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

// a system dependent on testSystem
type testDependentSystem struct {
	ts *testSystem `sameriver-system-dependency:"-"`
}

func newTestDependentSystem() *testDependentSystem {
	return &testDependentSystem{}
}
func (s *testDependentSystem) LinkWorld(w *World)   {}
func (s *testDependentSystem) Update(dt_ms float64) {}

// a system (misconfig) which is implemented on a non-pointer receiver
type testNonPointerReceiverSystem struct {
}

func newTestNonPointerReceiverSystem() testNonPointerReceiverSystem {
	return testNonPointerReceiverSystem{}
}
func (s testNonPointerReceiverSystem) LinkWorld(w *World)   {}
func (s testNonPointerReceiverSystem) Update(dt_ms float64) {}

// a system (misconfig) whose name does not end in System
type testSystemThatIsMisnamed struct {
}

func newTestSystemThatIsMisnamed() *testSystemThatIsMisnamed {
	return &testSystemThatIsMisnamed{}
}
func (s *testSystemThatIsMisnamed) LinkWorld(w *World)   {}
func (s *testSystemThatIsMisnamed) Update(dt_ms float64) {}

// a system (misconfig) which is dependent on a non-pointer type
type testDependentNonPointerSystem struct {
	ts testNonPointerReceiverSystem `sameriver-system-dependency:"-"`
}

func newTestDependentNonPointerSystem() *testDependentNonPointerSystem {
	return &testDependentNonPointerSystem{}
}
func (s *testDependentNonPointerSystem) LinkWorld(w *World)   {}
func (s *testDependentNonPointerSystem) Update(dt_ms float64) {}
