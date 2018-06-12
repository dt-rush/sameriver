package main

import (
	"math"
	"math/rand"
	"time"
)

type WorldMap struct {
	Lakes    []*Lake
	Vertices map[int]*MapVertex
	seed     int64
}

func (wm *WorldMap) NewVertex(p Point2D) *MapVertex {
	id := len(wm.Vertices)
	v := MapVertex{
		id:        id,
		pos:       p,
		neighbors: make([]*MapVertex, 0)}
	wm.Vertices[id] = &v
	return &v
}

func GenerateWorldMap() *WorldMap {
	seed := time.Now().UnixNano()
	m := WorldMap{seed: seed}
	m.Vertices = make(map[int]*MapVertex)
	nLakes := 5
	for i := 0; i < nLakes; i++ {
		m.AddRandomLake(Point2D{
			int((0.6*rand.Float64() + 0.2) * WORLD_WIDTH),
			int(rand.Float64() * WORLD_HEIGHT)})
	}
	return &m
}

func (wm *WorldMap) AddRandomLake(p Point2D) {

	l := Lake{}
	nVertices := 12 + rand.Intn(8)
	thetas := make([]float64, nVertices)
	jitter := 2 * math.Pi / float64(nVertices*2)
	smoothness := 0.4
	radiusBase := WORLD_WIDTH / (10.0 + 8*rand.Float64())
	var lastRadius = radiusBase
	for i := 0; i < nVertices; i++ {
		smoothness = smoothness * math.Cos(2*math.Pi*float64(i)/float64(nVertices))
		thetas[i] = float64(i) * (2 * math.Pi / float64(nVertices))
		thetas[i] += jitter * (rand.Float64() - 0.5)
		interp := (smoothness + (1-smoothness)*rand.Float64())
		radius := (1-interp)*lastRadius + interp*radiusBase*(3*(rand.Float64())+rand.Float64())
		lastRadius = radius
		vp := Point2D{
			p.X + int(radius*math.Cos(thetas[i])),
			p.Y + int(radius*math.Sin(thetas[i]))}
		v := wm.NewVertex(vp)
		l.Vertices = append(l.Vertices, v)
	}
	wm.Lakes = append(wm.Lakes, &l)
}
