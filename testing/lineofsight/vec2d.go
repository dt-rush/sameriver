package main

type Vec2D struct {
	X float64
	Y float64
}

func (v1 Vec2D) Sub(v2 Vec2D) Vec2D {
	return Vec2D{v1.X - v2.X, v1.Y - v2.Y}
}

func (v1 Vec2D) ScalarCross(v2 Vec2D) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}
