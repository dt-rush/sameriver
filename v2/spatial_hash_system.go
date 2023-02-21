package sameriver

// used to compute the spatial hash tables given a list of entities
type SpatialHashSystem struct {
	gridX  int
	gridY  int
	hasher *SpatialHasher
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
	s.hasher = NewSpatialHasher(s.gridX, s.gridY, w)
}

func (h *SpatialHashSystem) Update(dt_ms float64) {
	// clear any old data and run the computation
	h.hasher.ClearTable()
	h.hasher.ScanAndInsertEntities()
}

func (h *SpatialHashSystem) Expand(n int) {
	h.hasher.Expand(n)
}
