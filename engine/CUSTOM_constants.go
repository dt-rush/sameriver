package main

import (
	"time"
)

const WINDOW_TITLE = "SameRiver"

const SKIP_MENU = false

const WINDOW_WIDTH int32 = 600
const WINDOW_HEIGHT int32 = 450

const WORLD_WIDTH = 5760
const WORLD_HEIGHT = 5760

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond
