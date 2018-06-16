package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (g *Game) DrawWorld() {
	g.r.SetDrawColor(0, 0, 0, 255)
	g.r.FillRect(nil)
	g.DrawObstacles()
	g.DrawEntityAndPath()
	g.DrawDiffusionMap()
}

func (g *Game) DrawEntityAndPath() {

	if g.w.e != nil {
		if g.w.e.moveTarget != nil {
			drawPoint(g.r, *g.w.e.moveTarget,
				sdl.Color{R: 0, G: 255, B: 255}, POINTSZ)
			// toward := VecFromPoints(g.w.e.pos, *g.w.e.moveTarget).
			// Unit().Scale(VECLENGTH)
			// drawVector(g.r, g.w.e.pos, toward, sdl.Color{R: 255, G: 255, B: 0})
			drawVector(g.r, g.w.e.pos, g.w.e.vel.Scale(VECLENGTH), sdl.Color{R: 255, G: 0, B: 0})
		}
		drawPoint(g.r, g.w.e.pos,
			sdl.Color{R: 0, G: 255, B: 0}, POINTSZ)
	}
}

func (g *Game) DrawObstacles() {
	for _, o := range g.w.obstacles {
		drawRect(g.r, o, sdl.Color{R: 255, G: 0, B: 0})
	}
}

func (g *Game) DrawDiffusionMap() {
	g.r.Copy(g.w.dif.st, nil, nil)
}

func (g *Game) DrawUI() {
	g.r.Copy(g.ui.st, nil, nil)
}
