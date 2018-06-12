package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

func drawLake(r *sdl.Renderer, l *Lake) {
	vx := make([]int16, len(l.Vertices))
	vy := make([]int16, len(l.Vertices))
	for i, v := range l.Vertices {
		ssv := worldSpaceToScreenSpace(v.pos)
		vx[i] = int16(ssv.X)
		vy[i] = int16(ssv.Y)
	}
	gfx.FilledPolygonColor(r, vx, vy, sdl.Color{0, 0, 164, 255})
}

func (w *World) DrawWorldMap(r *sdl.Renderer) {
	r.SetDrawColor(0, 128, 0, 255)
	r.FillRect(&sdl.Rect{0, 0, WINDOW_WIDTH, WINDOW_HEIGHT})
	for _, l := range w.m.Lakes {
		drawLake(r, l)
	}
}

func drawPath(r *sdl.Renderer, p []Point2D) {

}

func drawPoint(r *sdl.Renderer, p *Point2D, c sdl.Color) {
	ssp := worldSpaceToScreenSpace(*p)
	r.SetDrawColor(c.R, c.G, c.B, 255)
	r.FillRect(&sdl.Rect{
		int32(ssp.X - 4),
		int32(ssp.Y - 4),
		8, 8})
}

func (w *World) DrawEntityAndPath(r *sdl.Renderer) {

	if w.e != nil {
		if w.e.path != nil {
			drawPath(r, w.e.path)
		}
		if w.e.moveTarget != nil {
			drawPoint(r, w.e.moveTarget, sdl.Color{R: 0, G: 255, B: 255})
		}
		drawPoint(r, &w.e.pos, sdl.Color{R: 255, G: 0, B: 0})
	}
}
