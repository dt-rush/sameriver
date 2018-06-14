package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"sync"
	"time"
)

type DiffusionMap struct {
	// values of the diffusion field
	d [DIFFUSION_DIM][DIFFUSION_DIM]float64
	// used in calculating the diffusion map
	inv      [DIFFUSION_DIM][DIFFUSION_DIM]bool
	invDirty sync.Mutex
	// ticker used to time updates to the map (since it can be expensive)
	tick *time.Ticker
	// array of obstacles in the world
	obstacles *[]Rect2D
	// renderer reference
	r *sdl.Renderer
	// screen texture
	st *sdl.Texture
}

func NewDiffusionMap(r *sdl.Renderer, obstacles *[]Rect2D) *DiffusionMap {
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
		tick:      time.NewTicker(100 * time.Millisecond),
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
			drawRect(m.r,
				Rect2D{
					float64(x) * DIFFUSION_CELL_W,
					float64(y) * DIFFUSION_CELL_H,
					DIFFUSION_CELL_W,
					DIFFUSION_CELL_H},
				sdl.Color{R: val, G: val, B: val})
		}
	}
}

func (m *DiffusionMap) Diffuse(pos Point2D) {
	m.invDirty.Lock()

	N := DIFFUSION_DIM * DIFFUSION_DIM

	initX := int(pos.X / DIFFUSION_CELL_W)
	initY := int(pos.Y / DIFFUSION_CELL_H)
	if initX > DIFFUSION_DIM-1 {
		initX = DIFFUSION_DIM - 1
	}
	if initY > DIFFUSION_DIM-1 {
		initY = DIFFUSION_DIM - 1
	}

	var validNeighbor = func(x int, y int) bool {
		return x >= 0 && x < DIFFUSION_DIM &&
			y >= 0 && y < DIFFUSION_DIM
	}

	var touch = func(x int, y int, d float64) {
		m.d[y][x] = d
	}

	frontier := make(map[[2]int]bool)
	frontier[[2]int{initX, initY}] = true
	touch(initX, initY, 1.0)
	for N > 0 {
		nextFrontier := make(map[[2]int]bool)
		for cur, _ := range frontier {
			// get diffusion field value for the popped frontier element
			d := m.d[cur[1]][cur[0]]
			// add valid neighbors to frontier, marking their val according to
			for iy := -1; iy <= 1; iy++ {
				for ix := -1; ix <= 1; ix++ {
					// a cell can't be it's own neighbor
					if ix == 0 && iy == 0 {
						continue
					}
					x := cur[0] + ix
					y := cur[1] + iy
					_, inFrontier := frontier[[2]int{x, y}]
					if !inFrontier && validNeighbor(x, y) {
						distance := d * math.Pow(0.98, math.Sqrt(float64(ix*ix+iy*iy)))
						if !m.inv[y][x] {
							m.inv[y][x] = true
							N--
							nextFrontier[[2]int{x, y}] = true
							touch(x, y, distance)
						} else if distance < m.d[y][x] {
							touch(x, y, distance)
						}
					}
				}
			}
		}
		for pos, _ := range frontier {
			delete(frontier, pos)
		}
		for pos, _ := range nextFrontier {
			frontier[pos] = true
		}
	}

	m.UpdateTexture()

	go func() {
		for y := 0; y < DIFFUSION_DIM; y++ {
			for x := 0; x < DIFFUSION_DIM; x++ {
				m.inv[y][x] = false
			}
		}
		m.invDirty.Unlock()
	}()
}

func Subtract() {
}
