/**
  * 
  * 
  * 
  * 
**/



// Game event constants

package engine

import (
	"github.com/dt-rush/donkeys-qquest/constants"
)


// actual game event constants defined in constants/constants.go


type GameEvent int

func (e GameEvent) String() string {
	// playing off the parallel array structure and
	// the fact that GameEvent is int
	return constants.GAME_EVENT_STRINGS [e]
}
