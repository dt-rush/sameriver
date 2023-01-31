package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

// mockup loading scene
type testingLoadingScene struct {
	initRan                bool
	updateRan              bool
	drawRan                bool
	handleKeyboardStateRan bool
	handleKeyboardEventRan bool
	nextSceneRan           bool
}

func (s *testingLoadingScene) Name() string {
	return "testingLoadingScene"
}
func (s *testingLoadingScene) Init(game *Game, config map[string]string) {
	s.initRan = true
}
func (s *testingLoadingScene) Update(dt_ms float64, allowance_ms float64) {
	s.updateRan = true
}
func (s *testingLoadingScene) Draw(window *sdl.Window, renderer *sdl.Renderer) {
	s.drawRan = true
}
func (s *testingLoadingScene) HandleKeyboardState(keyboard_state []uint8) {
	s.handleKeyboardStateRan = true
}
func (s *testingLoadingScene) HandleKeyboardEvent(keyboard_event *sdl.KeyboardEvent) {
	s.handleKeyboardEventRan = true
}
func (s *testingLoadingScene) IsDone() bool {
	return false
}
func (s *testingLoadingScene) NextScene() Scene {
	s.nextSceneRan = true
	return nil
}
func (s *testingLoadingScene) End() {
}
func (s *testingLoadingScene) IsTransient() bool {
	return false
}
func (s *testingLoadingScene) Destroy() {}
