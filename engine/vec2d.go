package engine

import (
	"math"
	"math/rand"
)

type Vec2D struct {
	X float32
	Y float32
}

func RandomUnitVec2D() Vec2D {
	x := rand.Float64()
	y := rand.Float64()
	l := math.Sqrt(x*x + y*y)
	x /= l
	y /= l
	return Vec2D{float32(x), float32(y)}

}
