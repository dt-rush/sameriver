package engine

import (
	"fmt"
	"testing"
)

func TestRectWithinRect(t *testing.T) {
	within := []RectPair{
		RectPair{
			Vec2D{4, 4},
			Vec2D{4, 4},
			Vec2D{4, 4},
			Vec2D{5, 5},
		},
		RectPair{
			Vec2D{0, 0},
			Vec2D{10, 10},
			Vec2D{0, 0},
			Vec2D{10, 10},
		},
		RectPair{
			Vec2D{0, 4},
			Vec2D{10, 6},
			Vec2D{4, 4},
			Vec2D{20, 20},
		},
		RectPair{
			Vec2D{6, 0},
			Vec2D{4, 10},
			Vec2D{6, 0},
			Vec2D{10, 10},
		},
	}
	notWithin := []RectPair{
		RectPair{
			Vec2D{0, 0},
			Vec2D{10, 10},
			Vec2D{4, 4},
			Vec2D{2, 2},
		},
		RectPair{
			Vec2D{10, 10},
			Vec2D{1, 1},
			Vec2D{0, 0},
			Vec2D{2, 2},
		},
	}

	for _, pair := range within {
		if !RectWithinRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should be within %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
	}
	for _, pair := range notWithin {
		if RectWithinRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should not be within %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
	}
}

func TestRectIntersectsRect(t *testing.T) {
	intersects := []RectPair{
		RectPair{
			Vec2D{4, 4},
			Vec2D{4, 4},
			Vec2D{0, 0},
			Vec2D{10, 10},
		},
		RectPair{
			Vec2D{2, 2},
			Vec2D{10, 10},
			Vec2D{0, 0},
			Vec2D{10, 10},
		},
	}
	doesntIntersect := []RectPair{
		RectPair{
			Vec2D{0, 0},
			Vec2D{4, 4},
			Vec2D{6, 6},
			Vec2D{4, 4},
		},
	}

	for _, pair := range intersects {
		if !RectIntersectsRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should intersect %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
		// swap rects and test again
		pair.pos0, pair.pos1 = pair.pos1, pair.pos0
		pair.box0, pair.box1 = pair.box1, pair.box0
		if !RectIntersectsRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should intersect %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
	}
	for _, pair := range doesntIntersect {
		if RectIntersectsRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should not intersect %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
		// swap rects and test again
		pair.pos0, pair.pos1 = pair.pos1, pair.pos0
		pair.box0, pair.box1 = pair.box1, pair.box0
		if RectIntersectsRect(&pair.pos0, &pair.box0, &pair.pos1, &pair.box1) {
			t.Fatal(fmt.Sprintf("%v,%v should not intersect %v,%v",
				&pair.pos0, &pair.box0, &pair.pos1, &pair.box1))
		}
	}
}
