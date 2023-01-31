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

const VERSION = "0.5.01"
const AUDIO_ON = true

const FPS = 60
const FRAME_DURATION = (1000 / FPS) * time.Millisecond
const FRAME_DURATION_INT = (1000 / FPS)
const MAX_ENTITIES = 1600

const COLLISION_RATELIMIT_TIMEOUT_MS = 300

const EVENT_PUBLISH_CHANNEL_CAPACITY = MAX_ENTITIES / 4
const EVENT_SUBSCRIBER_CHANNEL_CAPACITY = 64

const ADD_REMOVE_LOGIC_CHANNEL_CAPACITY = MAX_ENTITIES / 4
