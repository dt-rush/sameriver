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

func (s *SpatialHashSystem) GetComponentDeps() []string {
	return []string{"Vec2D,Position", "Vec2D,Box"}
}

func (s *SpatialHashSystem) LinkWorld(w *World) {
	s.Hasher = NewSpatialHasher(s.gridX, s.gridY, w)
}

func (s *SpatialHashSystem) Update(dt_ms float64) {
	// clear any old data and run the computation
	s.Hasher.ClearTable()
	s.Hasher.ScanAndInsertEntities()
}

func (s *SpatialHashSystem) Expand(n int) {
	s.Hasher.Expand(n)
}
