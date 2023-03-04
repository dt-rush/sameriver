package sameriver

import (
	"math"
)

// takes rectanges defined with pos in the center of the rect
func RectWithinRect(pos0, box0, pos1, box1 Vec2D) bool {
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
func RectIntersectsRect(pos0, box0, pos1, box1 Vec2D) bool {
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

func RectWithinRadiusOfPoint(pos, box Vec2D, d float64, point Vec2D) bool {
	// algorithm by MultiRRomero on stackoverflow
	// https://stackoverflow.com/a/18157551
	rectMin := pos.ShiftedCenterToBottomLeft(box)
	rectMax := rectMin.Add(box)
	dx := math.Max(rectMin.X-point.X, math.Max(0, point.X-rectMax.X))
	dy := math.Max(rectMin.Y-point.Y, math.Max(0, point.Y-rectMax.Y))
	return dx*dx+dy*dy < d*d
}

func RectWithinDistanceOfRect(iPos, iBox, jPos, jBox Vec2D, d float64) bool {
	dist := func(iPos, iBox, jPos, jBox Vec2D) (d float64) {
		x1 := iPos.X
		y1 := iPos.Y
		x1b := iPos.X + iBox.X
		y1b := iPos.Y + iBox.Y

		x2 := jPos.X
		y2 := jPos.Y
		x2b := jPos.X + jBox.X
		y2b := jPos.Y + jBox.Y

		// adapted from Maxim's stackoverflow answer
		// https://stackoverflow.com/a/26178015

		left := x2b < x1
		right := x1b < x2
		bottom := y2b < y1
		top := y1b < y2

		if top && left {
			_, _, d = Vec2D{x1, y1b}.Distance(Vec2D{x2b, y2})
		} else if left && bottom {
			_, _, d = Vec2D{x1, y1}.Distance(Vec2D{x2b, y2b})
		} else if bottom && right {
			_, _, d = Vec2D{x1b, y1}.Distance(Vec2D{x2, y2b})
		} else if right && top {
			_, _, d = Vec2D{x1b, y1b}.Distance(Vec2D{x2, y2})
		} else if left {
			return x1 - x2b
		} else if right {
			return x2 - x1b
		} else if bottom {
			return y1 - y2b
		} else if top {
			return y2 - y1b
		}
		return d
	}
	return dist(iPos, iBox, jPos, jBox) < d
}

/* TODO: reconsider this (even though it doesn't work right) if the above ever breaks?
func RectWithinDistanceOfRect(iPos, iBox, jPos, jBox Vec2D, d float64) bool {
	// algorithm from stackoverflow user Nick Alger
	// https://stackoverflow.com/a/65107290
	boxDist := func(aMin, aMax, bMin, bMax Vec2D) float64 {
		entrywiseMaxZero := func(vec Vec2D) Vec2D {
			return Vec2D{
				math.Max(0, vec.X),
				math.Max(0, vec.Y),
			}
		}
		euclidNorm := func(vec Vec2D) float64 {
			return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
		}
		u := entrywiseMaxZero(aMin.Sub(bMax))
		v := entrywiseMaxZero(bMin.Sub(aMax))
		unorm := euclidNorm(u)
		vnorm := euclidNorm(v)
		return math.Sqrt(unorm*unorm + vnorm*vnorm)
	}

	// lower-left and upper-right corners
	iMin := iPos.ShiftedCenterToBottomLeft(iBox)
	iMax := iMin.Add(iBox)
	jMin := jPos.ShiftedCenterToBottomLeft(jBox)
	jMax := jMin.Add(jBox)
	Logger.Printf("boxDist: %f", boxDist(iMin, iMax, jMin, jMax))
	return boxDist(iMin, iMax, jMin, jMax) < d
}
*/
