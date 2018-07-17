package engine

import (
	"regexp"
	"testing"
)

func TestEntityManagerConstruct(t *testing.T) {
	w := testingWorld()
	em := NewEntityManager(w)
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}

func TestEntityManagerNumEntities(t *testing.T) {
	w := testingWorld()
	total, active := w.em.NumEntities()
	if !(total == 0 && active == 0) {
		t.Fatal("somehow total and active were nonzero")
	}
	e, _ := w.em.spawnFromRequest(simpleSpawnRequest())
	total, active = w.em.NumEntities()
	if !(total == 1 && active == 1) {
		t.Fatal("total or active not updated after one spawn")
	}
	w.em.Deactivate(e)
	total, active = w.em.NumEntities()
	if !(total == 1 && active == 0) {
		t.Fatal("active count not updated or deactivate despawned")
	}
	w.em.doDespawn(e)
	total, active = w.em.NumEntities()
	if !(total == 0 && active == 0) {
		t.Fatal("did not despawn properly")
	}
}

func TestEntityManagerSpawnDespawn(t *testing.T) {
	w := testingWorld()
	w.em.spawnFromRequest(simpleSpawnRequest())
	total, _ := w.em.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	e := w.em.entities[0]
	if !e.Active {
		t.Fatal("Spawned entity was not active")
	}
	w.em.doDespawn(e)
	if e.Active {
		t.Fatal("deSpawn did not deactivate entity")
	}
	if !e.Despawned {
		t.Fatal("deSpawn did not set DeSpawned flag")
	}
}

func TestEntityManagerSpawnFail(t *testing.T) {
	w := testingWorld()

	for i := 0; i < MAX_ENTITIES; i++ {
		w.em.entityTable.allocateID()
	}
	_, err := w.em.spawnFromRequest(simpleSpawnRequest())
	if err == nil {
		t.Fatal("should have thrown error on spawnrequest when entity table full")
	}
}

func TestEntityManagerSpawnEvent(t *testing.T) {
	w := testingWorld()

	w.em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, simpleSpawnRequest()}
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ := w.em.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn subscription")
	}
}

func TestEntityManagerDespawnEvent(t *testing.T) {
	w := testingWorld()

	e, _ := w.em.spawnFromRequest(simpleSpawnRequest())
	w.em.despawnSubscription.C <- Event{
		DESPAWNREQUEST_EVENT,
		DespawnRequestData{Entity: e}}
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ := w.em.NumEntities()
	if total != 0 {
		t.Fatal("should have despawned an entity after processing despawn subscription")
	}
}

func TestEntityManagerDespawnAll(t *testing.T) {
	w := testingWorld()

	for i := 0; i < 64; i++ {
		w.em.spawnFromRequest(simpleSpawnRequest())
	}
	w.em.spawnSubscription.C <- Event{}
	w.em.despawnSubscription.C <- Event{}
	w.em.despawnAll()
	total, _ := w.em.NumEntities()
	if total != 0 {
		t.Fatal("did not despawn all entities")
	}
	if len(w.em.spawnSubscription.C) != 0 {
		t.Fatal("did not drain spawnSubscription channel")
	}
	if len(w.em.despawnSubscription.C) != 0 {
		t.Fatal("did not drain despawnSubscription channel")
	}
}

func TestEntityManagerEntityHasComponent(t *testing.T) {
	w := testingWorld()

	pos := Vec2D{11, 11}
	e, _ := w.em.spawnFromRequest(positionSpawnRequest(pos))
	if w.em.components.Position[e.ID] != pos {
		t.Fatal("failed to apply component data")
	}
	if !w.em.EntityHasComponent(e, POSITION_COMPONENT) {
		t.Fatal("failed to set entity component bit array")
	}
}

func TestEntityManagerEntitiesWithTag(t *testing.T) {
	w := testingWorld()

	tag := "tag1"
	req := simpleTaggedSpawnRequest(tag)
	w.em.spawnFromRequest(req)
	tagged := w.em.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	w := testingWorld()

	tag := "the chosen one"
	e, err := w.em.UniqueTaggedEntity(tag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	req := simpleSpawnRequest()
	_, err = w.em.spawnUnique(tag, req.Components, req.Tags)
	if err != nil {
		t.Fatal("failed to Spawn FIRST unique entity")
	}
	e, err = w.em.UniqueTaggedEntity(tag)
	if !(e != nil && err == nil) {
		t.Fatal("did not return unique entity")
	}
	_, err = w.em.spawnUnique(tag, req.Components, req.Tags)
	if err == nil {
		t.Fatal("should not have been allowed to Spawn second unique entity")
	}
}

func TestEntityManagerTagUntagEntity(t *testing.T) {
	w := testingWorld()

	w.em.spawnFromRequest(simpleSpawnRequest())
	e := w.em.entities[0]
	tag := "tag1"
	w.em.TagEntity(e, tag)
	tagged := w.em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
	if !w.em.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as untagged")
	}
	w.em.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
	if w.em.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as still having removed tag")
	}
}

func TestEntityManagerTagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e, _ := w.em.spawnFromRequest(simpleSpawnRequest())
		entities = append(entities, e)
	}
	w.em.TagEntities(entities, tag)
	for _, e := range entities {
		if !w.em.components.TagList[e.ID].Has(tag) {
			t.Fatal("entity's taglist was not modified by TagEntities")
		}
	}
}

func TestEntityManagerUntagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e, _ := w.em.spawnFromRequest(simpleTaggedSpawnRequest(tag))
		entities = append(entities, e)
	}
	w.em.UntagEntities(entities, tag)
	for _, e := range entities {
		if w.em.components.TagList[e.ID].Has(tag) {
			t.Fatal("entity's taglist was not modified by UntagEntities")
		}
	}
}

func TestEntityManagerDeactivateActivateEntity(t *testing.T) {
	w := testingWorld()

	w.em.spawnFromRequest(simpleSpawnRequest())
	e := w.em.entities[0]
	tag := "tag1"
	w.em.TagEntity(e, tag)
	tagged := w.em.EntitiesWithTag(tag)
	w.em.Deactivate(e)
	if tagged.Length() != 0 {
		t.Fatal("entity was not removed from list after Deactivate()")
	}
	_, active := w.em.NumEntities()
	if active != 0 {
		t.Fatal("didn't update active count")
	}
	w.em.Activate(e)
	if tagged.Length() != 1 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
	_, active = w.em.NumEntities()
	if active != 1 {
		t.Fatal("didn't update active count")
	}
}

func TestEntityManagerGetUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	name := "ILoveLily"
	list := w.em.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return true
		}),
	)

	w.em.spawnFromRequest(simpleSpawnRequest())
	if list.Length() != 1 {
		t.Fatal("failed to update UpdatedEntityList")
	}
	list2 := w.em.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return true
		}),
	)
	if list2.Length() != 1 {
		t.Fatal("failed to created UpdatedEntityList relative to existing entities")
	}
}

func TestEntityManagerGetSortedUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	list := w.em.GetSortedUpdatedEntityList(NewEntityFilter("filter",
		func(e *Entity) bool {
			return true
		}),
	)
	e8 := &Entity{ID: 8, Active: true, Despawned: false}
	e0 := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e8})
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	if list.entities[0].ID != 0 {
		t.Fatal("didn't insert in order")
	}
}

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	w := testingWorld()
	name := "ILoveLily"
	if w.em.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	list := w.em.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return false
		}),
	)
	if w.em.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func TestEntityManagerGetCurrentEntities(t *testing.T) {
	w := testingWorld()
	if !(len(w.em.GetCurrentEntities()) == 0) {
		t.Fatal("initially, len(GetCurrentEntities()) should be 0")
	}
	e, _ := w.em.spawnFromRequest(simpleSpawnRequest())
	if !(len(w.em.GetCurrentEntities()) == 1) {
		t.Fatal("after spawn, len(GetCurrentEntities()) should be 1")
	}
	w.em.Deactivate(e)
	if !(len(w.em.GetCurrentEntities()) == 1) {
		t.Fatal("after deactivate, len(GetCurrentEntities()) should be 1")
	}
	w.em.doDespawn(e)
	if !(len(w.em.GetCurrentEntities()) == 0) {
		t.Fatal("after despawn, len(GetCurrentEntities()) should be 0")
	}
}

func TestEntityManagerString(t *testing.T) {
	w := testingWorld()
	w.em.String()
}

func TestEntityManagerDump(t *testing.T) {
	w := testingWorld()
	w.em.spawnFromRequest(simpleSpawnRequest())
	e := w.em.entities[0]
	tag := "tag1"
	w.em.TagEntity(e, tag)
	s := w.em.Dump()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in EntityManager.Dump()")
	}
}
