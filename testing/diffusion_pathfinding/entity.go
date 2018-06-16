package main

import (
	"time"
)

type Entity struct {
	w *World

	pos        Vec2D
	vel        Vec2D
	moveTarget *Vec2D

	moveTicker *time.Ticker
	velTicker  *time.Ticker

	path []Vec2D
}

func NewEntity(pos Vec2D, w *World) *Entity {
	return &Entity{
		w:          w,
		pos:        pos,
		moveTicker: time.NewTicker(16 * time.Millisecond),
		velTicker:  time.NewTicker(42 * time.Millisecond),
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
		d := toward.Magnitude()
		force := toward.Unit().Scale(1)
		for _, o := range e.w.obstacles {
			// check if entity is moving toward colliding with obstacle
			for i := 0; i < 3; i++ {
				future := e.pos.Add(e.vel.Scale(float64(10 * (i + 1))))
				if o.Contains(future) {
					ovec := VecFromPoints(e.pos, Vec2D{o.X + o.W/2, o.Y + o.H/2})
					d := ovec.Magnitude()
					lv := ovec.PerpendicularUnit().Scale(1)
					if e.vel.Project(lv) > 0 {
						force = force.Add(lv.Scale(2 * e.vel.Magnitude() * sigma4(d, 36)))
					}
					rv := ovec.PerpendicularUnit().Scale(-1)
					if e.vel.Project(rv) > 0 {
						force = force.Add(rv.Scale(2 * e.vel.Magnitude() * sigma4(d, 36)))
					}
				}
			}
		}
		max := MOVESPEED * (1 - sigma4(d, 16))
		e.vel = e.vel.Add(force).Truncate(max)
	}
}

func (e *Entity) Move() {
	if e != nil {
		e.pos.X += e.vel.X
		e.pos.Y += e.vel.Y
		if e.moveTarget != nil {
			_, _, d := e.pos.Distance(*e.moveTarget)
			if d < MOVESPEED*2 {
				e.pos = *e.moveTarget
				e.vel = Vec2D{0, 0}
				e.moveTarget = nil
				return
			}
		}
	}
}
