package engine

import (
	"regexp"
	"testing"
)

func TestEntityManagerConstruct(t *testing.T) {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}

func TestEntityManagerSpawn(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.em.Spawn(simpleSpawnRequestData())
	total, _ := w.em.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	e := w.em.Entities[0]
	if !e.Active {
		t.Fatal("Spawned entity was not active")
	}
	w.em.Despawn(e)
	if e.Active {
		t.Fatal("deSpawn did not deactivate entity")
	}
	if !e.Despawned {
		t.Fatal("deSpawn did not set DeSpawned flag")
	}
}

func TestEntityManagerSpawnFail(t *testing.T) {
	w := NewWorld(1024, 1024)

	for i := 0; i < MAX_ENTITIES; i++ {
		w.em.entityTable.allocateID()
	}
	_, err := w.em.Spawn(simpleSpawnRequestData())
	if err == nil {
		t.Fatal("should have thrown error on spawnrequest when entity table full")
	}
	w.em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, simpleSpawnRequestData()}
	w.em.Update()
}

func TestEntityManagerSpawnRequest(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, simpleSpawnRequestData()}
	w.Update(FRAME_SLEEP_MS)
	total, _ := w.em.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn subscription")
	}
}

func TestEntityManagerDespawnRequest(t *testing.T) {
	w := NewWorld(1024, 1024)

	e, _ := w.em.Spawn(simpleSpawnRequestData())
	w.em.despawnSubscription.C <- Event{
		DESPAWNREQUEST_EVENT,
		DespawnRequestData{Entity: e}}
	w.Update(FRAME_SLEEP_MS)
	total, _ := w.em.NumEntities()
	if total != 0 {
		t.Fatal("should have despawned an entity after processing despawn subscription")
	}
}

func TestEntityManagerDespawnAll(t *testing.T) {
	w := NewWorld(1024, 1024)

	for i := 0; i < 64; i++ {
		w.em.Spawn(simpleSpawnRequestData())
	}
	w.em.spawnSubscription.C <- Event{}
	w.em.despawnSubscription.C <- Event{}
	w.em.DespawnAll()
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

func TestEntityManagerSpawnWithComponent(t *testing.T) {
	w := NewWorld(1024, 1024)

	pos := Vec2D{11, 11}
	e, _ := w.em.Spawn(positionSpawnRequestData(pos))
	if w.em.Components.Position[e.ID] != pos {
		t.Fatal("failed to apply component data")
	}
	if !w.em.EntityHasComponent(e, POSITION_COMPONENT) {
		t.Fatal("failed to set entity component bit array")
	}
}

func TestEntityManagerEntitiesWithTagList(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "tag1"
	req := simpleTaggedSpawnRequestData(tag)
	w.em.Spawn(req)
	w.em.Update()
	tagged := w.em.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "the chosen one"
	e, err := w.em.UniqueTaggedEntity(tag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	req := simpleSpawnRequestData()
	_, err = w.em.SpawnUnique(tag, req)
	if err != nil {
		t.Fatal("failed to Spawn FIRST unique entity")
	}
	e, err = w.em.UniqueTaggedEntity(tag)
	if !(e != nil && err == nil) {
		t.Fatal("did not return unique entity")
	}
	w.em.Update()
	_, err = w.em.SpawnUnique(tag, req)
	if err == nil {
		t.Fatal("should not have been allowed to Spawn second unique entity")
	}
}

func TestEntityManagerTagUntagEntity(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.em.Spawn(simpleSpawnRequestData())
	w.em.Update()
	e := w.em.Entities[0]
	tag := "tag1"
	w.em.TagEntity(e, tag)
	tagged := w.em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
	w.em.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
}

func TestEntityManagerTagUntagEntities(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.em.Spawn(simpleSpawnRequestData())
	w.em.Spawn(simpleSpawnRequestData())
	w.em.Spawn(simpleSpawnRequestData())
	w.em.Update()
	tag := "tag1"
	w.em.TagEntities(w.em.GetCurrentEntities(), tag)
	tagged := w.em.EntitiesWithTag(tag)
	if tagged.Length() != 3 {
		t.Fatal("failed to tag all entities")
	}
	w.em.UntagEntities(w.em.GetCurrentEntities(), tag)
	if tagged.Length() != 0 {
		t.Fatal("failed to untag all entities")
	}
}

func TestEntityManagerDeactivateActivateEntity(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.em.Spawn(simpleSpawnRequestData())
	w.em.Update()
	e := w.em.Entities[0]
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

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	w := NewWorld(1024, 1024)
	name := "ILoveLily"
	if w.em.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	query := EntityQuery{
		Name: name,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			return false
		}}
	list := w.em.GetUpdatedEntityList(query)
	if w.em.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func TestEntityManagerString(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.em.String()
}

func TestEntityManagerDump(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.em.Spawn(simpleSpawnRequestData())
	w.em.Update()
	e := w.em.Entities[0]
	tag := "tag1"
	w.em.TagEntity(e, tag)
	s := w.em.Dump()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in EntityManager.Dump()")
	}
}
