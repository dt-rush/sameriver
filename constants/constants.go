/**
  * 
  * This file defines GAME constants 
  * 
  * 
**/



package constants

import (
	"fmt"
)



var VERSION = "0.1.0"
var WINDOW_TITLE = fmt.Sprintf ("Donkeys QQuest")
var WINDOW_WIDTH int32 = 400
var WINDOW_HEIGHT int32 = 300
var DEBUG_COLLISION = true



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
