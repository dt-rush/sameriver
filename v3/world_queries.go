package sameriver

func (w *World) PredicateAllEntities(p func(*Entity) bool) []*Entity {
	entities := make([]*Entity, 0)
	for e := range w.GetCurrentEntitiesSet() {
		entities = append(entities, e)
	}
	return w.PredicateEntities(entities, p)
}

func (w *World) PredicateEntities(entities []*Entity, p func(*Entity) bool) []*Entity {
	results := make([]*Entity, 0)
	for _, e := range entities {
		if !e.Active {
			continue
		}
		if p(e) {
			results = append(results, e)
		}
	}
	return results
}

func (w *World) EntitiesWithTags(tags ...string) []*Entity {
	entities := make([]*Entity, 0)
	for e := range w.GetCurrentEntitiesSet() {
		has := true
		for _, tag := range tags {
			has = has && w.EntityHasTag(e, tag)
		}
		if has {
			entities = append(entities, e)
		}
	}
	return entities
}

func (w *World) ActiveEntitiesWithTags(tags ...string) []*Entity {
	entities := make([]*Entity, 0)
	for e := range w.GetCurrentEntitiesSet() {
		if !e.Active {
			continue
		}
		has := true
		for _, tag := range tags {
			has = has && w.EntityHasTag(e, tag)
		}
		if has {
			entities = append(entities, e)
		}
	}
	return entities
}

func (w *World) EntitiesWithinDistance(pos, box Vec2D, d float64) []*Entity {
	return w.SpatialHasher.EntitiesWithinDistance(pos, box, d)
}

func (w *World) EntitiesWithinDistanceFilter(
	pos, box Vec2D, d float64, filter func(*Entity) bool) []*Entity {
	return w.SpatialHasher.EntitiesWithinDistanceFilter(pos, box, d, filter)
}

func (w *World) EntitiesWithinDistanceApprox(pos, box Vec2D, d float64) []*Entity {
	return w.SpatialHasher.EntitiesWithinDistanceApprox(pos, box, d)
}

func (w *World) EntitiesWithinDistanceApproxFilter(
	pos, box Vec2D, d float64, filter func(*Entity) bool) []*Entity {
	return w.SpatialHasher.EntitiesWithinDistanceApproxFilter(pos, box, d, filter)
}

func (w *World) CellsWithinDistance(pos, box Vec2D, d float64) [][2]int {
	return w.SpatialHasher.CellsWithinDistance(pos, box, d)
}

func (w *World) CellsWithinDistanceApprox(pos, box Vec2D, d float64) [][2]int {
	return w.SpatialHasher.CellsWithinDistanceApprox(pos, box, d)
}
