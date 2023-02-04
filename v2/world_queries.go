package sameriver

func (w *World) PredicateEntities(p func(*Entity) bool) []*Entity {
	results := make([]*Entity, 0)
	for e, _ := range w.em.entityTable.currentEntities {
		if !e.Active {
			continue
		}
		if p(e) {
			results = append(results, e)
		}
	}
	return results
}

/*
func (w *World) CellsWithinDistance(e *Entity, d float64) [][2]int {
	if sh, ok := w.systems["SpatialHashSystem"]; !ok {
		panic("Can't call CellsWithinDistance without SpatialHashSystem!")
	} else {
		return sh.(*SpatialHashSystem).CellsWithinDistance(*e.GetVec2D("Position"), d)
	}
}

func (w *World) EntitiesPotentiallyWithinDistance(e *Entity, d float64) []*Entity {
	if sh, ok := w.systems["SpatialHashSystem"]; !ok {
		panic("Can't call EntitiesPotentaillyWithinDistance without SpatialHashSystem!")
	} else {
		return sh.(*SpatialHashSystem).EntitiesPotentiallyWithinDistance(e, d)
	}
}

func (w *World) EntitiesWithinDistance(e *Entity, d float64) []*Entity {
	if sh, ok := w.systems["SpatialHashSystem"]; !ok {
		panic("Can't call EntitiesWithinDistance without SpatialHashSystem!")
	} else {
		return sh.(*SpatialHashSystem).EntitiesWithinDistance(e, d)
	}
}
*/
