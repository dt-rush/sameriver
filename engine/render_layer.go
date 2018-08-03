package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type RenderFunc func(w *sdl.Window, r *sdl.Renderer)

type RenderLayer struct {
	renderFunc RenderFunc
	z          int
	active     bool
}

func NewRenderLayer(renderFunc RenderFunc, z int) *RenderLayer {
	return &RenderLayer{
		renderFunc: renderFunc,
		z:          z,
		active:     true,
	}
}

func (l *RenderLayer) Activate() {
	l.active = true
}

func (l *RenderLayer) Deactivate() {
	l.active = false
}

func (l *RenderLayer) IsActive() bool {
	return l.active
}

func (l *RenderLayer) Render(w *sdl.Window, r *sdl.Renderer) {
	l.renderFunc(w, r)
}
