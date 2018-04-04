/**
  *
  *
  *
  *
**/

package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
)

type AudioManager struct {
	audio map[string](*mix.Chunk)
}

func (m *AudioManager) Play(file string) {
	if m.audio[file] == nil {
		return
	}
	// play on channel 1 (so that sounds cut each other off)
	// loop 0 times
	m.audio[file].Play(1, 0)
}

func (m *AudioManager) Load(file string) {
	log_err := func(err error) {
		Logger.Printf("failed to load assets/%s", file)
		panic(err)
	}
	chunk, err := mix.LoadWAV(fmt.Sprintf("assets/%s", file))
	if err != nil {
		log_err(err)
		m.audio[file] = nil
	} else {
		m.audio[file] = chunk
	}
}

func (m *AudioManager) Init() {
	// can be tuned
	capacity := 4
	m.audio = make(map[string](*mix.Chunk), capacity)
}
