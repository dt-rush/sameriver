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

const VERSION = "0.4.0"
const AUDIO_ON = true

const DEBUG_ENTITY_MANAGER = false
const DEBUG_EVENTS = false
const DEBUG_UPDATED_ENTITY_LISTS = false
const DEBUG_GOROUTINES = false
const DEBUG_ENTITY_LOGIC = false
const DEBUG_ENTITY_MANAGER_UPDATE_TIMING = false
const DEBUG_SPAWN = false
const DEBUG_DESPAWN = false
const DEBUG_ATOMIC_MODIFY = false
const DEBUG_ENTITY_CLASS = false
const DEBUG_WORLD_LOGIC = false
const DEBUG_ENTITY_LOCKS = false
const DEBUG_BEHAVIOR = false
const DEBUG_TAGS = false
const DEBUG_FUNCTION_END = false

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond
const MAX_ENTITIES = 1600

const COLLISION_RATELIMIT_TIMEOUT_MS = 500

const SPAWN_CHANNEL_CAPACITY = 128
const ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY = MAX_ENTITIES / 4
const EVENT_PUBLISH_CHANNEL_CAPACITY = MAX_ENTITIES / 4
const EVENT_SUBSCRIBER_CHANNEL_CAPACITY = 16
