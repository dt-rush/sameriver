package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
)

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
		g.Init(WindowSpec{
			Title:      "testing game",
			Width:      100,
			Height:     100,
			Fullscreen: false},
			&gameScene)
		g.Run()
		if !(expectedLoadingScene.initRan == loadingScene.initRan &&
			expectedLoadingScene.updateRan == loadingScene.updateRan &&
			expectedLoadingScene.drawRan == loadingScene.drawRan &&
			expectedLoadingScene.handleKeyboardStateRan == loadingScene.handleKeyboardStateRan &&
			expectedLoadingScene.handleKeyboardEventRan == loadingScene.handleKeyboardEventRan &&
			expectedLoadingScene.nextSceneRan == loadingScene.nextSceneRan) {
			t.Fatal("pattern of method calls did not match expected for loadingscene")
		}
		if !(expectedGameScene.initRan == gameScene.initRan &&
			expectedGameScene.updateRan == gameScene.updateRan &&
			expectedGameScene.drawRan == gameScene.drawRan &&
			expectedGameScene.handleKeyboardStateRan == gameScene.handleKeyboardStateRan &&
			expectedGameScene.handleKeyboardEventRan == gameScene.handleKeyboardEventRan &&
			expectedGameScene.nextSceneRan == gameScene.nextSceneRan) {
			t.Fatal("pattern of method calls did not match expected for loadingscene")
		}
	})
}
