/**
  * 
  * 
  * 
  * 
**/



package systems


import (

	"github.com/veandco/go-sdl2/sdl"

	// TODO  re: big philosophical debate lower down
	// figure if these are needed, it might be that all
	// this system needs to do is spawn entities with components,
	// leaving the textures, etc. up to
	// the renderer of those components
	
//	"github.com/veandco/go-sdl2/ttf"
//	"github.com/veandco/go-sdl2/img"

)

// responsible for spawning screen message entities
// managing their lifecycles, and destroying their resources
// when needed
type ScreenmessageSystem struct {
	// TODO  re: big philosophical debate lower down
	// if messages is map of string -> int
	// are those ints ID's of the message entities, or
	// are they indexes into an array of textures?
	// surely the textures need positions, etc.
	// starts to seem that these should be entities, yeah
	messages map[string]int
	textures []*sdl.Texture
}

// TODO big philoosophical debate here: SHOULD MENUS AND MESSAGES BE ENTITIES TOO?

func (s *ScreenmessageSystem) init (capacity int) {
	s.messages = make (map[string]int, capacity)
	s.textures = make ([]*sdl.Texture, capacity)
}
