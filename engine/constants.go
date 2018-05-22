/**
  *
  * This file defines constants for the game engine operation,
  * mostly flags to be set at compile-time.
  *
  *
**/

package engine

const VERSION = "0.3.2"
const AUDIO_ON = true

const DEBUG_GAME_EVENTS = true
const DEBUG_UPDATED_ENTITY_LISTS = false
const DEBUG_GOROUTINES = true

const FPS = 60
const MAX_ENTITIES = 1024

const COLLISION_RATELIMIT_TIMEOUT_MS = 500

const SPAWN_CHANNEL_CAPACITY = 128
const COMPONENT_MODIFICATION_CHANNEL_CAPACITY = MAX_ENTITIES
const ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY = 16
const EVENT_CHANNEL_CAPACITY = 16
