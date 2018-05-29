package main

import (
	"math/rand"
	"time"
)

const N_ENTITIES = 2000
const N_ENTITIES_WITH_BEHAVIOR = 1024
const N_BEHAVIORS_PER_ENTITY = 8
const CHANCE_TO_SAFEGET_IN_BEHAVIOR = 0.3
const SAFEGET_DURATION = 4 * time.Microsecond
const CHANCE_TO_LOCK_OTHER_ENTITY = 0.05
const OTHER_ENTITY_LOCK_DURATION = 12 * time.Microsecond

const WORLD_WIDTH = 3200
const WORLD_HEIGHT = 3200
const GRID = 32
const CELL_WIDTH = WORLD_WIDTH / GRID
const CELL_HEIGHT = WORLD_HEIGHT / GRID
const UPPER_ESTIMATE_ENTITIES_PER_SQUARE = 24
const GRID_GOROUTINES = 12

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond

var ATOMIC_MODIFY_DURATION = func() time.Duration {
	return time.Duration(rand.Intn(4000)) * time.Nanosecond
}

var BEHAVIOR_POST_SLEEP = func() time.Duration {
	return time.Duration(time.Duration(rand.Intn(700)) * time.Millisecond)
}

const PRINT_EXAMPLE = false
