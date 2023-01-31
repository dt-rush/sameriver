package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

// mockup game scene
type testingGameScene struct {
	accum_ms float64

	initRan                bool
	updateRan              bool
	drawRan                bool
	handleKeyboardStateRan bool
	handleKeyboardEventRan bool
	nextSceneRan           bool
}

func (s *testingGameScene) Name() string {
	return "testingGameScene"
}
func (s *testingGameScene) Init(game *Game, config map[string]string) {
	s.initRan = true
	time.Sleep(8 * FRAME_DURATION)
}
func (s *testingGameScene) Update(dt_ms float64, allowance_ms float64) {
	s.updateRan = true
	s.accum_ms += dt_ms
}
func (s *testingGameScene) Draw(window *sdl.Window, renderer *sdl.Renderer) {
	s.drawRan = true
}
func (s *testingGameScene) HandleKeyboardState(keyboard_state []uint8) {
	s.handleKeyboardStateRan = true
}
func (s *testingGameScene) HandleKeyboardEvent(keyboard_event *sdl.KeyboardEvent) {
	s.handleKeyboardEventRan = true
}
func (s *testingGameScene) IsDone() bool {
	return s.accum_ms >= 8*FRAME_DURATION_INT
}
func (s *testingGameScene) NextScene() Scene {
	s.nextSceneRan = true
	return nil
}
func (s *testingGameScene) End() {
}
func (s *testingGameScene) IsTransient() bool {
	return false
}
func (s *testingGameScene) Destroy() {}
