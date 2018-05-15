package engine

type CollisionLogic struct {
	Selector func(id_a uint16,
		id_b uint16,
		em *EntityManager) bool
	EventGenerator func(id_a uint16,
		id_b uint16,
		em *EntityManager) GameEvent
}
