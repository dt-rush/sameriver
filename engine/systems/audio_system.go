/**
  *
  *
  *
  *
**/



package systems

import (
    "github.com/veandco/go-sdl2/mix"
    "github.com/dt-rush/donkeys-qquest/engine"
)

type AudioSystem struct {
    audio map[string](*mix.Chunk)
}

func (s *AudioSystem) Play (file string) {
    if s.audio[file] == nil {
        return
    }
    // play on the first unreserved channel (-1), and
    // loop 0 times
    s.audio[file].Play (1, 0)
}

func (s *AudioSystem) Load (file string) {
    log_err := func (err error) {
        engine.Logger.Printf ("failed to load %s", file)
        panic (err)
    }
    chunk, err := mix.LoadWAV (file)
    if err != nil {
        log_err (err)
        s.audio [file] = nil
    } else {
        s.audio [file] = chunk
    }
}

func (s *AudioSystem) Init (capacity int) {
    s.audio = make (map[string](*mix.Chunk), capacity)
}
