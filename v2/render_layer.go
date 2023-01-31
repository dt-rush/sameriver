package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

type RenderFunc func(w *sdl.Window, r *sdl.Renderer)

type RenderLayer struct {
	name       string
	z          int
	renderFunc RenderFunc
	active     bool
}

func NewRenderLayer(name string, z int, renderFunc RenderFunc) *RenderLayer {
	return &RenderLayer{
		name:       name,
		z:          z,
		renderFunc: renderFunc,
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

func (l *RenderLayer) Name() string {
	return l.name
}

func (l *RenderLayer) Render(w *sdl.Window, r *sdl.Renderer) {
	l.renderFunc(w, r)
}
