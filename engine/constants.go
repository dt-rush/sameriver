/**
  *
  * This file defines constants for the game engine operation,
  * mostly flags to be set at compile-time.
  *
  *
**/

package engine

import (
	"time"
)

const VERSION = "0.5.01"
const AUDIO_ON = true

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond
const FRAME_SLEEP_MS = (1000 / FPS)
const MAX_ENTITIES = 1600

const COLLISION_RATELIMIT_TIMEOUT_MS = 300

const EVENT_PUBLISH_CHANNEL_CAPACITY = MAX_ENTITIES / 4
const EVENT_SUBSCRIBER_CHANNEL_CAPACITY = 64
