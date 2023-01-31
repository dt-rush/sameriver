package sameriver

type RectPair struct {
	pos0 Vec2D
	box0 Vec2D
	pos1 Vec2D
	box1 Vec2D
}

// takes rectanges defined with pos in the center of the rect
func RectWithinRect(pos0, box0, pos1, box1 *Vec2D) bool {
	pos0.ShiftCenterToBottomLeft(box0)
	pos1.ShiftCenterToBottomLeft(box1)
	defer pos0.ShiftBottomLeftToCenter(box0)
	defer pos1.ShiftBottomLeftToCenter(box1)
	return pos0.X >= pos1.X &&
		pos0.X+box0.X <= pos1.X+box1.X &&
		pos0.Y >= pos1.Y &&
		pos0.Y+box0.Y <= pos1.Y+box1.Y
}

// translated from SDL2's SDL_HasIntersection(SDL_Rect * A, SDL_Rect * B)
// takes rectanges defined with pos in the center of the rect
func RectIntersectsRect(pos0, box0, pos1, box1 *Vec2D) bool {
	pos0.ShiftCenterToBottomLeft(box0)
	pos1.ShiftCenterToBottomLeft(box1)
	defer pos0.ShiftBottomLeftToCenter(box0)
	defer pos1.ShiftBottomLeftToCenter(box1)
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
