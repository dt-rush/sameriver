package sameriver

func (w *World) PredicateAllEntities(p func(*Entity) bool) []*Entity {
	entities := make([]*Entity, 0)
	for e, _ := range w.em.entityTable.currentEntities {
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
	for e, _ := range w.em.entityTable.currentEntities {
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
	for e, _ := range w.em.entityTable.currentEntities {
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
