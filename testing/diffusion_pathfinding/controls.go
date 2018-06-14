package main

type ControlMode int

const N_MODES = 2

const (
	MODE_PLACING_WAYPOINT = 0
	MODE_PLACING_OBSTACLE = iota
)

var MODENAMES []string = []string{
	"MODE_PLACING_WAYPOINT",
	"MODE_PLACING_OBSTACLE",
}

type Controls struct {
	mode ControlMode
}

func NewControls() *Controls {
	return &Controls{mode: MODE_PLACING_WAYPOINT}
}

func (c *Controls) ToggleMode() {
	c.mode = (c.mode + 1) % N_MODES
}
