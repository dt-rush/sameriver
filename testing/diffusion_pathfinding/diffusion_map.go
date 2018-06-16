package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type DiffusionMap struct {
	// values of the diffusion field
	d [DIFFUSION_DIM][DIFFUSION_DIM]float64
	// ticker used to time updates to the map (since it can be expensive)
	tick *time.Ticker
	// array of obstacles in the world
	obstacles *[]Rect2D
	// renderer reference
	r *sdl.Renderer
	// screen texture
	st *sdl.Texture
}

func NewDiffusionMap(
	r *sdl.Renderer, obstacles *[]Rect2D, tick time.Duration) *DiffusionMap {
	st, err := r.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)
	st.SetBlendMode(sdl.BLENDMODE_BLEND)
	st.SetAlphaMod(0x80)
	if err != nil {
		panic(err)
	}

	return &DiffusionMap{
		obstacles: obstacles,
		tick:      time.NewTicker(tick),
		r:         r,
		st:        st}
}

func (m *DiffusionMap) UpdateTexture() {
	m.r.SetRenderTarget(m.st)
	defer m.r.SetRenderTarget(nil)

	m.r.SetDrawColor(0, 0, 0, 0)
	m.r.Clear()
	for y := 0; y < DIFFUSION_DIM; y++ {
		for x := 0; x < DIFFUSION_DIM; x++ {
			val := uint8(255 * m.d[y][x])
			var c sdl.Color
			if m.CellHasObstacle(x, y) {
				c = sdl.Color{R: val, G: 0, B: 0}
			} else {
				c = sdl.Color{R: val, G: val, B: val}
			}
			drawRect(m.r,
				Rect2D{
					float64(x) * DIFFUSION_CELL_W,
					float64(y) * DIFFUSION_CELL_H,
					DIFFUSION_CELL_W,
					DIFFUSION_CELL_H},
				c,
			)
		}
	}
}

func (m *DiffusionMap) CellHasObstacle(x int, y int) bool {
	r := Rect2D{
		float64(x) * DIFFUSION_CELL_W,
		float64(y) * DIFFUSION_CELL_H,
		DIFFUSION_CELL_W,
		DIFFUSION_CELL_H}
	for _, o := range *m.obstacles {
		if r.Overlaps(o) {
			return true
		}
	}
	return false
}

func (m *DiffusionMap) Diffuse(pos Vec2D) {
	initX := int(pos.X / DIFFUSION_CELL_W)
	initY := int(pos.Y / DIFFUSION_CELL_H)
	if initX > DIFFUSION_DIM-1 {
		initX = DIFFUSION_DIM - 1
	}
	if initY > DIFFUSION_DIM-1 {
		initY = DIFFUSION_DIM - 1
	}
	m.d[initY][initX] = 1.0

	var validNeighbor = func(x int, y int) bool {
		return x >= 0 && x < DIFFUSION_DIM &&
			y >= 0 && y < DIFFUSION_DIM
	}

	var avgOfNeighbors = func(x int, y int) float64 {
		neumann := [][2]int{
			[2]int{-1, 0},
			[2]int{1, 0},
			[2]int{0, -1},
			[2]int{0, 1},
		}
		sum := 0.0
		n := 0.0
		for _, neu := range neumann {
			ox := x + neu[0]
			oy := y + neu[1]
			if validNeighbor(ox, oy) {
				sum += m.d[oy][ox]
				n++
			}
		}
		return sum / n
	}

	for y := 0; y < DIFFUSION_DIM; y++ {
		for x := 0; x < DIFFUSION_DIM; x++ {
			if y == initY && x == initX {
				continue
			}
			m.d[y][x] = avgOfNeighbors(x, y)
			m.d[y][x] *= 0.998
		}
	}

	m.UpdateTexture()
}

func Subtract() {
}
