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

func (h *SpatialHashSystem) Update(dt_ms float64) {
	// clear any old data and run the computation
	h.Hasher.ClearTable()
	h.Hasher.ScanAndInsertEntities()
}

func (h *SpatialHashSystem) Expand(n int) {
	h.Hasher.Expand(n)
}
