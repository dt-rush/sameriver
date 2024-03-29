package sameriver

import (
	"time"
)

// a basic system with a data member
type testSystem struct {
	updates int
}

func newTestSystem() *testSystem {
	return &testSystem{}
}
func (s *testSystem) LinkWorld(w *World) {}
func (s *testSystem) Update(dt_ms float64) {
	time.Sleep(1 * time.Millisecond)
	s.updates += 1
}
func (s *testSystem) GetComponentDeps() []any {
	return []any{}
}
func (s *testSystem) Expand(n int) {}

// a system dependent on testSystem
type testDependentSystem struct {
	ts *testSystem `sameriver-system-dependency:"-"`
}

func newTestDependentSystem() *testDependentSystem {
	return &testDependentSystem{}
}
func (s *testDependentSystem) LinkWorld(w *World)   {}
func (s *testDependentSystem) Update(dt_ms float64) {}
func (s *testDependentSystem) GetComponentDeps() []any {
	return []any{}
}
func (s *testDependentSystem) Expand(n int) {}

// a system (misconfig) which is implemented on a non-pointer receiver
type testNonPointerReceiverSystem struct {
}

func newTestNonPointerReceiverSystem() testNonPointerReceiverSystem {
	return testNonPointerReceiverSystem{}
}
func (s testNonPointerReceiverSystem) LinkWorld(w *World)   {}
func (s testNonPointerReceiverSystem) Update(dt_ms float64) {}
func (s testNonPointerReceiverSystem) GetComponentDeps() []any {
	return []any{}
}
func (s testNonPointerReceiverSystem) Expand(n int) {}

// a system (misconfig) whose name does not end in System
type testSystemThatIsMisnamed struct {
}

func newTestSystemThatIsMisnamed() *testSystemThatIsMisnamed {
	return &testSystemThatIsMisnamed{}
}
func (s *testSystemThatIsMisnamed) LinkWorld(w *World)   {}
func (s *testSystemThatIsMisnamed) Update(dt_ms float64) {}
func (s *testSystemThatIsMisnamed) GetComponentDeps() []any {
	return []any{}
}
func (s *testSystemThatIsMisnamed) Expand(n int) {}

// a system (misconfig) which is dependent on a non-pointer type
type testDependentNonPointerSystem struct {
	ts testNonPointerReceiverSystem `sameriver-system-dependency:"-"`
}

func newTestDependentNonPointerSystem() *testDependentNonPointerSystem {
	return &testDependentNonPointerSystem{}
}
func (s *testDependentNonPointerSystem) LinkWorld(w *World)   {}
func (s *testDependentNonPointerSystem) Update(dt_ms float64) {}
func (s *testDependentNonPointerSystem) GetComponentDeps() []any {
	return []any{}
}

func (s *testDependentNonPointerSystem) Expand(n int) {}

// a system (misconfig) which is dependent on a non-System
type testDependentNonSystemSystem struct {
	ts *EntityManager `sameriver-system-dependency:"-"`
}

func newTestDependentNonSystemSystem() *testDependentNonSystemSystem {
	return &testDependentNonSystemSystem{}
}
func (s *testDependentNonSystemSystem) LinkWorld(w *World)   {}
func (s *testDependentNonSystemSystem) Update(dt_ms float64) {}
func (s *testDependentNonSystemSystem) GetComponentDeps() []any {
	return []any{}
}
func (s *testDependentNonSystemSystem) Expand(n int) {}
