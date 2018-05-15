/**
  *
  * Manages the loading and playback of Audio resources
  *
  *
**/

package engine

import (
	"fmt"
	"io/ioutil"

	"github.com/veandco/go-sdl2/mix"
)

// AudioManager stores audio as mix.Chunk pointers,
// keyed by strings (filenames)
type AudioManager struct {
	audio map[string](*mix.Chunk)
}

// Init the map which stores the audio chunks
func (m *AudioManager) Init() {
	// can be tuned
	capacity := 4
	m.audio = make(map[string](*mix.Chunk), capacity)
	// read all audio files in assets/audio
	files, err := ioutil.ReadDir("assets/audio")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		m.Load(f.Name())
	}
}

// loads an audio file in the assets/ folder into the map, making it playable
func (m *AudioManager) Load(file string) {
	chunk, err := mix.LoadWAV(fmt.Sprintf("assets/audio/%s", file))
	if err != nil {
		Logger.Printf("[Audio manager] failed to load assets/%s", file)
		m.audio[file] = nil
	} else {
		m.audio[file] = chunk
	}
}

// on execution of this function, the given audio will begin to play
func (m *AudioManager) Play(file string) {
	if m.audio[file] == nil {
		// the value in the map will be nil if the asset
		// failed to load in Load()
		Logger.Printf("[Audio manager] attempted to play asset %s, which had failed to load",
			file)
		return
	} else {
		// play on channel 1 (so that sounds cut each other off)
		// loop 0 times
		m.audio[file].Play(1, 0)
	}
}
