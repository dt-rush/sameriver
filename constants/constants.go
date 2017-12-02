/**
  * 
  * This file defines GAME constants 
  * 
  * 
**/



package constants



/*
*             NOTE ON THE BELOW
*
*       FOR THE LOVE OF ALL THAT IS HOLY
*       KEEP THE CONSTANT DECLARATION STATEMENT AND THE ARRAY ALIGNED
*/

const (
	GAME_EVENT_DONKEY_CAUGHT = iota
	GAME_EVENT_FLAME_HIT_PLAYER = iota
)

// used to support String()
var GAME_EVENT_STRINGS = []string{
	"GAME_EVENT_DONKEY_CAUGHT",
	"GAME_EVENT_FLAME_HIT_PLAYER"}
