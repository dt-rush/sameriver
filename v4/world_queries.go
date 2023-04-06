package sameriver

import "math"

func (w *World) FilterAllEntities(filter func(*Entity) bool) []*Entity {
	results := make([]*Entity, 0)
	for e := range w.GetCurrentEntitiesSet() {
		if !e.Active {
			continue
		}
		if filter(e) {
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
			has = has && e.HasTag(tag)
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
			has = has && e.HasTag(tag)
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

func (w *World) ClosestEntityFilter(pos Vec2D, box Vec2D, filter func(*Entity) bool) *Entity {
	closest := (*Entity)(nil)
	closestDistance := math.MaxFloat64
	for _, e := range w.FilterAllEntities(filter) {
		entityPos := e.GetVec2D(POSITION)
		entityBox := e.GetVec2D(BOX)
		distance := RectDistance(pos, box, *entityPos, *entityBox)
		if distance < closestDistance {
			closestDistance = distance
			closest = e
		}
	}
	return closest
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
