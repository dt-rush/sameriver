package engine

type ComponentType int

const N_COMPONENT_TYPES = 7

const (
	ACTIVE_COMPONENT   = iota
	COLOR_COMPONENT    = iota
	HEALTH_COMPONENT   = iota
	HITBOX_COMPONENT   = iota
	LOGIC_COMPONENT    = iota
	MIND_COMPONENT     = iota
	POSITION_COMPONENT = iota
	SPRITE_COMPONENT   = iota
	TAGLIST_COMPONENT  = iota
	VELOCITY_COMPONENT = iota
)

var COMPONENT_NAMES = map[ComponentType]string{
	ACTIVE_COMPONENT:   "active component",
	COLOR_COMPONENT:    "color component",
	HEALTH_COMPONENT:   "health component",
	HITBOX_COMPONENT:   "hitbox component",
	LOGIC_COMPONENT:    "logic component",
	MIND_COMPONENT:     "mind component",
	POSITION_COMPONENT: "position component",
	SPRITE_COMPONENT:   "sprite component",
	TAGLIST_COMPONENT:  "taglist component",
	VELOCITY_COMPONENT: "velocity component",
}
