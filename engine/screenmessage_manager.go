/**
  *
  *
  *
  *
**/

package engine

type ScreenMessage interface {
	Position() [2]int
	Text() string
	Update(dt_ms int)
	IsActive() bool
}

type FixedScreenMessage struct {
	// the message to display
	Msg string
	// how many lines of text to show before needing the
	// player to press space
	Lines int
	// age in milliseconds (used to scroll the text (if applicable)
	// and to set inactive)
	Age int
}

type FloatingScreenMessage struct {
	// the message to display
	Msg string
	// the top-left corner of the box, where (0, 0) is
	// the bottom-left corner of the screen
	Position [2]int
	// how long the message should float for (in milliseconds)
	Duration int
	// used to time the disappearance of the message
	Age int
}

// responsible for spawning screen message entities
// managing their lifecycles, and destroying their resources
// when needed
type ScreenMessageManager struct {
	messages map[int]ScreenMessage
}

func (s *ScreenMessageManager) Init() {
	// arbitrary, can be tuned? Will grow
	capacity := 4
	s.messages = make(map[int]ScreenMessage, capacity)
}

func (s *ScreenMessageManager) Update(dt_ms int) {
	// TODO: implement
}
