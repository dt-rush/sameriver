package main

import (
	"time"
)

type Entity struct {
	w *World

	pos        Vec2D
	vel        Vec2D
	moveTarget *Vec2D
	path       []Vec2D

	moveTicker *time.Ticker
	velTicker  *time.Ticker
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
		var nextPathPoint *Vec2D = nil
		for nextPathPoint == nil {
			if len(e.path) == 0 {
				nextPathPoint = e.moveTarget
				break
			}
			nextPathPoint = &e.path[len(e.path)-1]
			if e.pos.Sub(*nextPathPoint).Magnitude() < GRIDCELL_WORLD_W {
				e.path = e.path[:len(e.path)-1]
				nextPathPoint = nil
				continue
			}
		}
		toward := (*nextPathPoint).Sub(e.pos)
		d := toward.Magnitude()
		force := toward.Unit()
		avoid := 4.0
		for _, o := range e.w.obstacles {
			// check if entity is moving toward colliding with obstacle
			oc := Vec2D{o.X + o.W/2, o.Y + o.H/2}
			ovec := oc.Sub(e.pos)
			if d > 4*OBSTACLESZ &&
				e.vel.Project(ovec) > 0 && ovec.Magnitude() < OBSTACLESZ*3 {
				lv := ovec.PerpendicularUnit().Scale(1)
				if e.vel.Project(lv) > 0 {
					force = force.Add(lv.Scale(avoid * OBSTACLESZ *
						sigma4(ovec.Magnitude(), OBSTACLESZ/2)))
				}
				rv := ovec.PerpendicularUnit().Scale(-1)
				if e.vel.Project(rv) > 0 {
					force = force.Add(rv.Scale(avoid * OBSTACLESZ *
						sigma4(ovec.Magnitude(), OBSTACLESZ/2)))
				}
			}
		}
		max := MOVESPEED * (1 - sigma4(e.moveTarget.Sub(e.pos).Magnitude(), 16))
		e.vel = e.vel.Add(force).Truncate(max)
	}
}

func (e *Entity) Move() {
	if e != nil {
		vel := e.vel
		vX := e.vel.XComponent()
		vY := e.vel.YComponent()
		for _, o := range e.w.obstacles {
			re := Rect2D{
				e.pos.X - ENTITYSZ/2 - 2, e.pos.Y - ENTITYSZ/2 - 2,
				ENTITYSZ + 4, ENTITYSZ + 4}
			reX := re.Add(vX)
			if o.Overlaps(reX) {
				dxL := o.X - (e.pos.X + ENTITYSZ/2)
				dxR := e.pos.X - (o.X + o.W + ENTITYSZ/2)
				if dxL > 0 {
					// we are to the left
					vX = vX.Truncate(dxL * 0.5)
				} else if dxR > 0 {
					// we are to the right
					vX = vX.Truncate(dxR * 0.5)
				}
			}
			if o.Overlaps(re.Add(vX)) {
				vX = Vec2D{0, 0}
			}
			reY := re.Add(vY)
			if o.Overlaps(reY) {
				dyD := o.Y - (e.pos.Y + ENTITYSZ/2)
				dyU := e.pos.Y - (o.Y + o.H + ENTITYSZ/2)
				if dyD > 0 {
					// we are down
					vY = vY.Truncate(dyD * 0.5)
				} else if dyU > 0 {
					// we are up
					vY = vY.Truncate(dyU * 0.5)
				}
			}
			if o.Overlaps(re.Add(vY)) {
				vY = Vec2D{0, 0}
			}
		}
		vel = vX.Add(vY)
		e.pos = e.pos.Add(vel)
	}
}
