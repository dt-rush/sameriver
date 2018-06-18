package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (g *Game) DrawGRID() {
	g.r.SetDrawColor(0, 0, 0, 255)
	g.r.FillRect(nil)
	if g.showData {
		g.DrawDiffusionMap()
	}
	g.DrawObstacles()
	g.DrawEntityAndPath()
}

func (g *Game) DrawEntityAndPath() {

	if g.w.e != nil {

		drawPoint(g.r, g.w.e.pos,
			sdl.Color{R: 0, G: 255, B: 0}, ENTITYSZ)

		if g.w.e.moveTarget != nil {

			drawPoint(g.r, *g.w.e.moveTarget,
				sdl.Color{R: 0, G: 255, B: 255}, POINTSZ)

			if g.showData {
				toward := g.w.e.moveTarget.Sub(g.w.e.pos).Unit().Scale(VECLENGTH)
				drawVector(g.r, g.w.e.pos, toward, sdl.Color{R: 255, G: 255, B: 0})
				drawVector(g.r, g.w.e.pos, g.w.e.vel.Scale(VECLENGTH), sdl.Color{R: 255, G: 0, B: 0})

				for _, p := range g.w.e.path {
					drawPoint(g.r, p, sdl.Color{R: 128, G: 180, B: 255}, 3)
				}

			}

		}

	}
}

func (g *Game) DrawObstacles() {
	for _, o := range g.w.obstacles {

		bodyColor := sdl.Color{R: 255, G: 0, B: 0}
		pointColor := sdl.Color{R: 0, G: 0, B: 0}
		if g.w.e != nil && g.w.e.moveTarget != nil {
			toward := g.w.e.moveTarget.Sub(g.w.e.pos)
			d := toward.Magnitude()
			oc := Vec2D{o.X + o.W/2, o.Y + o.H/2}
			ovec := oc.Sub(g.w.e.pos)
			if d > 4*OBSTACLESZ &&
				g.w.e.vel.Project(ovec) > 0 && ovec.Magnitude() < OBSTACLESZ*3 {
				bodyColor = sdl.Color{R: 128, G: 0, B: 0}
				lv := ovec.PerpendicularUnit().Scale(1)
				rv := ovec.PerpendicularUnit().Scale(-1)
				if g.w.e.vel.Project(lv) > 0 || g.w.e.vel.Project(rv) > 0 {
					pointColor = sdl.Color{R: 255, G: 255, B: 255}
				}
			}

		}
		drawRect(g.r, o, bodyColor)
		center := Vec2D{o.X + o.W/2, o.Y + o.H/2}
		drawPoint(g.r, center, pointColor, 3)

	}
}

func (g *Game) DrawDiffusionMap() {
	g.r.Copy(g.w.dm.st, nil, nil)
}

func (g *Game) DrawUI() {
	g.r.Copy(g.ui.st, nil, nil)
}
