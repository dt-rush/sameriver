package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
	"time"
)

//
// mockup loading scene
//
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
func (s *testingLoadingScene) Update(dt_ms float64) {
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
func (s *testingLoadingScene) IsTransient() bool {
	return false
}
func (s *testingLoadingScene) Destroy() {}

//
// mockup game scene
//
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
	time.Sleep(8 * FRAME_SLEEP)
}
func (s *testingGameScene) Update(dt_ms float64) {
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
	return s.accum_ms >= 8*FRAME_SLEEP_MS
}
func (s *testingGameScene) NextScene() Scene {
	s.nextSceneRan = true
	return nil
}
func (s *testingGameScene) IsTransient() bool {
	return false
}
func (s *testingGameScene) Destroy() {}

// test functions
func TestGameLoadingSceneGameScene(t *testing.T) {
	// skip if CI (needs video device)
	skipCI(t)
	sdl.Main(func() {
		expectedLoadingScene := testingLoadingScene{
			initRan:                true,
			updateRan:              true,
			drawRan:                true,
			handleKeyboardStateRan: true,
			handleKeyboardEventRan: false,
			nextSceneRan:           true,
		}
		expectedGameScene := testingGameScene{
			initRan:                true,
			updateRan:              true,
			drawRan:                true,
			handleKeyboardStateRan: true,
			handleKeyboardEventRan: false,
			nextSceneRan:           true,
		}
		g := NewGame()
		loadingScene := testingLoadingScene{}
		g.SetLoadingScene(&loadingScene)
		gameScene := testingGameScene{}
		g.Init("testing game", 100, 100, &gameScene)
		g.Run()
		if !(expectedLoadingScene.initRan == loadingScene.initRan &&
			expectedLoadingScene.updateRan == loadingScene.updateRan &&
			expectedLoadingScene.drawRan == loadingScene.drawRan &&
			expectedLoadingScene.handleKeyboardStateRan == loadingScene.handleKeyboardStateRan &&
			expectedLoadingScene.handleKeyboardEventRan == loadingScene.handleKeyboardEventRan &&
			expectedLoadingScene.nextSceneRan == loadingScene.nextSceneRan) {
			t.Fatal("pattern of method calls did not match expected for loadingscene")
		}
		if !(expectedGameScene.initRan == loadingScene.initRan &&
			expectedGameScene.updateRan == loadingScene.updateRan &&
			expectedGameScene.drawRan == loadingScene.drawRan &&
			expectedGameScene.handleKeyboardStateRan == loadingScene.handleKeyboardStateRan &&
			expectedGameScene.handleKeyboardEventRan == loadingScene.handleKeyboardEventRan &&
			expectedGameScene.nextSceneRan == loadingScene.nextSceneRan) {
			t.Fatal("pattern of method calls did not match expected for loadingscene")
		}
	})
}
