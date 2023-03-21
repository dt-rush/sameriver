/**
  *
  * This file defines constants for the game engine operation,
  * mostly flags to be set at compile-time.
  *
  *
**/

package sameriver

import (
	"time"
)

const AUDIO_ON = true

const FPS = 60
const FRAME_DURATION = (1000 / FPS) * time.Millisecond
const FRAME_MS = (1000 / FPS)
const MAX_ENTITIES = 4096

// a subscriber getting 4096 events in a single update tick is insane,
// but memory is plentiful so, allow some capacity to build up
const EVENT_SUBSCRIBER_CHANNEL_CAPACITY = 128

const ADD_REMOVE_LOGIC_CHANNEL_CAPACITY = MAX_ENTITIES / 4

const RUNTIME_LIMIT_SHARER_MAX_LOOPS = 8
