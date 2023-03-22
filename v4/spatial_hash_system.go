package sameriver

// used to compute the spatial hash tables given a list of entities
type SpatialHashSystem struct {
	gridX  int
	gridY  int
	Hasher *SpatialHasher
}

func NewSpatialHashSystem(gridX, gridY int) *SpatialHashSystem {
	return &SpatialHashSystem{
		gridX: gridX,
		gridY: gridY,
	}
}

func (s *SpatialHashSystem) GetComponentDeps() map[ComponentID]ComponentKind {
	return map[ComponentID]ComponentKind{
		POSITION: VEC2D,
		BOX:      VEC2D,
	}
}

func (s *SpatialHashSystem) LinkWorld(w *World) {
	s.Hasher = NewSpatialHasher(s.gridX, s.gridY, w)
}

func (s *SpatialHashSystem) Update(dt_ms float64) {
	// clear any old data and run the computation
	s.Hasher.Update()
}

func (s *SpatialHashSystem) Expand(n int) {
	s.Hasher.Expand(n)
}
