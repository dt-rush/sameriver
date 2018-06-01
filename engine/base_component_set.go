package engine

type BaseComponentSet struct {
	Position *[2]int16
	HitBox   *[2]uint16
	Sprite   *Sprite
	TagList  *TagList
	Velocity *[2]float32
}
