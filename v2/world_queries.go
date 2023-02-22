package sameriver

func (w *World) PredicateAllEntities(p func(*Entity) bool) []*Entity {
	entities := make([]*Entity, 0)
	for e, _ := range w.GetCurrentEntitiesSet() {
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
	for e, _ := range w.GetCurrentEntitiesSet() {
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
	for e, _ := range w.GetCurrentEntitiesSet() {
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
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call EntitiesWithinDistance without SpatialHashSystem registered")
	}
	return sh.Hasher.EntitiesWithinDistance(pos, box, d)
}

func (w *World) EntitiesWithinDistanceFilter(pos, box Vec2D, d float64, filter func(*Entity) bool) []*Entity {
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call EntitiesWithinDistanceFilter without SpatialHashSystem registered")
	}
	all := sh.Hasher.EntitiesWithinDistance(pos, box, d)
	result := make([]*Entity, 0)
	for _, e := range all {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

func (w *World) EntitiesWIthinDistanceApprox(pos, box Vec2D, d float64) []*Entity {
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call EntitiesWithinDistanceApprox without SpatialHashSystem registered")
	}
	return sh.Hasher.EntitiesWithinDistanceApprox(pos, box, d)
}

func (w *World) EntitiesWithinDistanceApproxFilter(pos, box Vec2D, d float64, filter func(*Entity) bool) []*Entity {
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call EntitiesWithinDistanceApproxFilter without SpatialHashSystem registered")
	}
	all := sh.Hasher.EntitiesWithinDistanceApprox(pos, box, d)
	result := make([]*Entity, 0)
	for _, e := range all {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

func (w *World) CellsWithinDistance(pos, box Vec2D, d float64) [][2]int {
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call CellsWithinDistance without SpatialHashSystem registered")
	}
	return sh.Hasher.CellsWithinDistance(pos, box, d)
}

func (w *World) CellsWithinDistanceApprox(pos, box Vec2D, d float64) [][2]int {
	sh, ok := w.systems["SpatialHashSystem"].(*SpatialHashSystem)
	if !ok {
		panic("Tried to call CellsWithinDistanceApprox without SpatialHashSystem registered")
	}
	return sh.Hasher.CellsWithinDistanceApprox(pos, box, d)
}
