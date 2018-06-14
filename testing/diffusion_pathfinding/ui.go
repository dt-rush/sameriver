package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type UI struct {
	// renderer
	r *sdl.Renderer
	// font
	f *ttf.Font
	// screen texture
	st *sdl.Texture
	// message
	msg string
}

func NewUI(r *sdl.Renderer, f *ttf.Font) *UI {
	st, err := r.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)
	st.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return &UI{r: r, f: f, st: st}
}

func (ui *UI) UpdateMsg(msg string) {
	ui.msg = msg
	ui.renderMsgToST(msg,
		sdl.Color{255, 255, 255, 255},
		&sdl.Rect{20, 20,
			int32(WINDOW_WIDTH / 80 * len(msg)),
			int32(WINDOW_HEIGHT / 40)})
}

// render message to screen texture
func (ui *UI) renderMsgToST(
	msg string,
	color sdl.Color,
	dst *sdl.Rect) {

	ui.r.SetRenderTarget(ui.st)
	defer ui.r.SetRenderTarget(nil)

	ui.r.SetDrawColor(0, 0, 0, 0)
	ui.r.Clear()

	// surface & texture to be used in writing the
	// message to the texture
	var surface *sdl.Surface
	var texture *sdl.Texture
	var err error
	surface, err = ui.f.RenderUTF8Solid(msg, color)
	if err != nil {
		panic(err)
	}
	texture, err = ui.r.CreateTextureFromSurface(
		surface)
	if err != nil {
		panic(err)
	}
	// this copies our texture to the *target* texture
	// (screen_texture)
	ui.r.SetDrawColor(
		color.R,
		color.G,
		color.B,
		color.A)
	ui.r.Copy(texture, nil, dst)
	// free the resources allocated above
	surface.Free()
	texture.Destroy()
}
