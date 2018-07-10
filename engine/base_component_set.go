package engine

type BaseComponentSet struct {
	TagList *TagList

	Logic *LogicUnit

	Box    *Vec2D
	Sprite *Sprite

	Mass           *float64
	Position       *Vec2D
	Velocity       *Vec2D
	MaxVelocity    *float64
	MovementTarget *Vec2D
	Steer          *Vec2D
}
