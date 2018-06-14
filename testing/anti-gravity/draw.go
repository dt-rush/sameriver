package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

func drawPoint(r *sdl.Renderer, p Point2D, c sdl.Color, sz int) {
	sspx, sspy := WorldSpaceToScreenSpace(p)
	r.SetDrawColor(c.R, c.G, c.B, 255)
	r.FillRect(&sdl.Rect{
		int32(int(sspx) - sz/2),
		int32(int(sspy) - sz/2),
		int32(sz), int32(sz)})
}

func drawVector(r *sdl.Renderer, pos Point2D, v Vec2D, c sdl.Color) {
	// screen-space position
	sspx, sspy := WorldSpaceToScreenSpace(pos)
	// screen-space vector tip
	ssvtx, ssvty := WorldSpaceToScreenSpace(
		PointDelta(pos, v.X, v.Y))
	gfx.LineColor(r,
		int32(sspx),
		int32(sspy),
		int32(ssvtx),
		int32(ssvty),
		sdl.Color{R: c.R, G: c.G, B: c.B, A: 255})
}

func (g *Game) DrawWorld() {
	g.r.SetDrawColor(0, 0, 0, 255)
	g.r.FillRect(nil)
	g.DrawObstacles()
	g.DrawEntityAndPath()
}

func (g *Game) DrawEntityAndPath() {

	if g.w.e != nil {
		if g.w.e.moveTarget != nil {
			drawPoint(g.r, *g.w.e.moveTarget,
				sdl.Color{R: 0, G: 255, B: 255}, POINTSZ)
			toward := VecFromPoints(g.w.e.pos, *g.w.e.moveTarget).
				Unit().Scale(VECLENGTH)
			drawVector(g.r, g.w.e.pos, toward, sdl.Color{R: 255, G: 255, B: 0})

			ov := VecFromPoints(g.w.e.pos, *g.w.e.moveTarget)
			pl := ov.PerpendicularUnit().Scale(-16)
			pr := ov.PerpendicularUnit().Scale(16)
			drawVector(g.r, g.w.e.pos, pl, sdl.Color{R: 0, G: 255, B: 255})
			drawVector(g.r, g.w.e.pos, pr, sdl.Color{R: 0, G: 0, B: 255})

		}
		drawPoint(g.r, g.w.e.pos,
			sdl.Color{R: 0, G: 255, B: 0}, POINTSZ)
	}
}

func (g *Game) DrawObstacles() {
	for _, o := range g.w.obstacles {
		drawPoint(g.r, o,
			sdl.Color{R: 255, G: 0, B: 0}, OBSTACLESZ)
	}
}

func (g *Game) DrawUI() {
	g.ui.mutex.Lock()
	defer g.ui.mutex.Unlock()
	g.r.Copy(g.ui.st, nil, nil)
}
