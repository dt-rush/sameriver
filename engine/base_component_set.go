package engine

type BaseComponentSet struct {
	TagList *TagList

	Logic *LogicUnit

	Box    *Vec2D
	Sprite *Sprite

	Position       *Vec2D
	Velocity       *Vec2D
	MovementTarget *Vec2D
	Steer          *float64
}
