package engine

type Vec2D struct {
	x float32
	y float32
}

func RandomUnitVec2D() Vec2D {
	x := rand.Float64()
	y := rand.Float64()
	l := math.Sqrt(x*x + y*y)
	x /= l
	y /= l
	return Vec2D{x, y}
}
