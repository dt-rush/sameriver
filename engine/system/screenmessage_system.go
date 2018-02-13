/**
  *
  *
  *
  *
**/

package system

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type ScreenMessage struct {
	msg      string
	floating bool
	// the top-left corner of the box, where (0, 0) is
	// the bottom-left corner of the screen
	position_x int
	position_y int
}

// responsible for spawning screen message entities
// managing their lifecycles, and destroying their resources
// when needed
type ScreenMessageSystem struct {
	messages   map[string]int
	textures   map[int]*sdl.Texture
	small_font *ttf.Font
}

func (s *ScreenMessageSystem) Init(capacity int) {
	var err error
	s.messages = make(map[string]int, capacity)
	s.textures = make(map[int]*sdl.Texture, capacity)
	if s.small_font, err = ttf.OpenFont("assets/fixedsys.ttf", 9); err != nil {
		panic(err)
	}
}

func (s *ScreenMessageSystem) DisplayScreenMessage(m ScreenMessage) {
	// render the texture

}
