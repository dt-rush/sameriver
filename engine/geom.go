package engine

type RectPair struct {
	pos0 Vec2D
	box0 Vec2D
	pos1 Vec2D
	box1 Vec2D
}

func RectWithinRect(pos0 Vec2D, box0 Vec2D, pos1 Vec2D, box1 Vec2D) bool {
	// NOTE: pos is bottom-left
	return pos0.X >= pos1.X &&
		pos0.X+box0.X <= pos1.X+box1.X &&
		pos0.Y >= pos1.Y &&
		pos0.Y+box0.Y <= pos1.Y+box1.Y
}

// translated from SDL2's SDL_HasIntersection(SDL_Rect * A, SDL_Rect * B)
func RectIntersectsRect(pos0 Vec2D, box0 Vec2D, pos1 Vec2D, box1 Vec2D) bool {
	// horizontal
	Amin := pos0.X
	Amax := Amin + box0.X
	Bmin := pos1.X
	Bmax := Bmin + box1.X
	if Bmin > Amin {
		Amin = Bmin
	}
	if Bmax < Amax {
		Amax = Bmax
	}
	if Amax <= Amin {
		return false
	}
	// vertical
	Amin = pos0.Y
	Amax = Amin + box0.Y
	Bmin = pos1.Y
	Bmax = Bmin + box1.Y
	if Bmin > Amin {
		Amin = Bmin
	}
	if Bmax < Amax {
		Amax = Bmax
	}
	if Amax <= Amin {
		return false
	}
	return true
}
