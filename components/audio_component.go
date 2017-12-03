/**
  * 
  * 
  * 
  * 
**/



package components

import (
	"fmt"

//	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/mix"
	"github.com/dt-rush/donkeys-qquest/engine"
)

type AudioComponent struct {
	// entity ID -> audio index
	data map[int](int)
	// audio index -> *mix.Music
	audio [](*mix.Music)
	// for getting the index into the *mix.Music array
	// via the original filename
	name_index_map map[string]int
}

func (c *AudioComponent) IndexOf (s string) int {
	return c.name_index_map [s]
}

func (c *AudioComponent) Get (id int) *mix.Music {
	return c.audio [c.data [id]]
}

func (c *AudioComponent) Set (id int, val interface{}) {
	val_ := val.(int)
	c.data[id] = val_
}

func (c *AudioComponent) DefaultValue () interface{} {
	return 0
}

func (c *AudioComponent) String() string {
	return fmt.Sprintf ("%v", c.data)
}

func (c *AudioComponent) Name() string {
	return "AudioComponent"
}

func (c *AudioComponent) Has (id int) bool {
	_, ok := c.data[id]
	return ok
}

func (c *AudioComponent) Init (capacity int, game *engine.Game) {

	// init data storage
	c.data = make (map[int]int, capacity)
	c.audio = make ([]*mix.Music, capacity)
	c.name_index_map = make (map[string]int, capacity)


	// audio file enum for now, dynamic load later (TODO)
	const (
		SUCCESS = "assets/success.wav"
	)
	to_load := []string{SUCCESS}
	for i, s := range to_load {
		var err error
		audio, err := mix.LoadMUS (s)
		if err != nil {
			fmt.Println (err)
			continue
		}
		fmt.Println ("about to write to ", i, " in c.audio")
		fmt.Println (audio)
		c.audio = append (c.audio, audio)
		fmt.Println ("finished?")
	}

}

// TODO implement
// func (c *AudioComponent) destroy() {

// }
