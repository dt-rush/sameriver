package engine

type CollisionQuery struct {
	Test func(id_a uint16,
		id_b uint16,
		em *EntityManager) bool
}
