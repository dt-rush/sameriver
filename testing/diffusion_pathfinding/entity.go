package main

import (
	"time"
)

type Entity struct {
	w *World

	pos        Point2D
	vel        Vec2D
	moveTarget *Point2D

	moveTicker *time.Ticker
	velTicker  *time.Ticker

	path []Point2D
}

func NewEntity(pos Point2D, w *World) *Entity {
	return &Entity{
		w:          w,
		pos:        pos,
		moveTicker: time.NewTicker(16 * time.Millisecond),
		velTicker:  time.NewTicker(96 * time.Millisecond),
	}
}

func (e *Entity) Update() {
	select {
	case _ = <-e.moveTicker.C:
		e.Move()
	case _ = <-e.velTicker.C:
		e.UpdateVel()
	default:
	}
}

func (e *Entity) UpdateVel() {

	if e.moveTarget != nil {
		toward := VecFromPoints(e.pos, *e.moveTarget)
		force := toward.Unit().
			Scale(1 + 16*MOVESPEED/toward.Magnitude())
		e.vel = e.vel.Add(force).Truncate(MOVESPEED)
	}
}

func (e *Entity) Move() {
	if e != nil {
		e.pos.X += e.vel.X
		e.pos.Y += e.vel.Y
		if e.moveTarget != nil {
			_, _, d := Distance(e.pos, *e.moveTarget)
			if d < MOVESPEED*2 {
				e.pos = *e.moveTarget
				e.vel = Vec2D{0, 0}
				e.moveTarget = nil
				return
			}
		}
	}
}
