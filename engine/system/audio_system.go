/**
  *
  *
  *
  *
**/

package system

import (
	"fmt"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/veandco/go-sdl2/mix"
)

type AudioSystem struct {
	audio map[string](*mix.Chunk)
}

func (s *AudioSystem) Play(file string) {
	if s.audio[file] == nil {
		return
	}
	// play on channel 1 (so that sounds cut each other off)
	// loop 0 times
	s.audio[file].Play(1, 0)
}

func (s *AudioSystem) Load(file string) {
	log_err := func(err error) {
		engine.Logger.Printf("failed to load assets/%s", file)
		panic(err)
	}
	chunk, err := mix.LoadWAV(fmt.Sprintf("assets/%s", file))
	if err != nil {
		log_err(err)
		s.audio[file] = nil
	} else {
		s.audio[file] = chunk
	}
}

func (s *AudioSystem) Init(capacity int) {
	s.audio = make(map[string](*mix.Chunk), capacity)
}
