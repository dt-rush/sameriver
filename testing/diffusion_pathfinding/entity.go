package main

import (
	"fmt"
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

func sigma(x float64, m float64) float64 {
	return 1.0 / ((x*x)/(m*m) + 1)
}

func (e *Entity) UpdateVel() {

	if e.moveTarget != nil {
		toward := VecFromPoints(e.pos, *e.moveTarget)
		d := toward.Magnitude()
		force := toward.Unit().Scale(1)
		for _, o := range e.w.obstacles {
			ovec := VecFromPoints(e.pos, Point2D{o.X + o.W/2, o.Y + o.H/2})
			d := ovec.Magnitude()
			lv := ovec.PerpendicularUnit().Scale(1)
			if e.vel.Project(lv) > 0 {
				fmt.Println("obstacle is to the left")
				force = force.Add(lv.Scale(8 * e.vel.Magnitude() * sigma(d, 16)))
			}
			rv := ovec.PerpendicularUnit().Scale(-1)
			if e.vel.Project(rv) > 0 {
				force = force.Add(rv.Scale(8 * e.vel.Magnitude() * sigma(d, 16)))
			}
		}
		max := MOVESPEED * (1 - sigma(d, 16))
		e.vel = e.vel.Add(force).Truncate(max)
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
