package engine

import (
	"regexp"
	"testing"
)

func TestEntityManagerConstruct(t *testing.T) {
	w := NewWorld(1024, 1024)
	em := NewEntityManager(w)
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}

func TestEntityManagerSpawn(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.Em.Spawn(simpleSpawnRequestData())
	total, _ := w.Em.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	e := w.Em.Entities[0]
	if !e.Active {
		t.Fatal("Spawned entity was not active")
	}
	w.Em.Despawn(e)
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
		w.Em.entityTable.allocateID()
	}
	_, err := w.Em.Spawn(simpleSpawnRequestData())
	if err == nil {
		t.Fatal("should have thrown error on spawnrequest when entity table full")
	}
	w.Em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, simpleSpawnRequestData()}
	w.Em.Update()
}

func TestEntityManagerSpawnRequest(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.Em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, simpleSpawnRequestData()}
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ := w.Em.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn subscription")
	}
}

func TestEntityManagerDespawnRequest(t *testing.T) {
	w := NewWorld(1024, 1024)

	e, _ := w.Em.Spawn(simpleSpawnRequestData())
	w.Em.despawnSubscription.C <- Event{
		DESPAWNREQUEST_EVENT,
		DespawnRequestData{Entity: e}}
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ := w.Em.NumEntities()
	if total != 0 {
		t.Fatal("should have despawned an entity after processing despawn subscription")
	}
}

func TestEntityManagerDespawnAll(t *testing.T) {
	w := NewWorld(1024, 1024)

	for i := 0; i < 64; i++ {
		w.Em.Spawn(simpleSpawnRequestData())
	}
	w.Em.spawnSubscription.C <- Event{}
	w.Em.despawnSubscription.C <- Event{}
	w.Em.DespawnAll()
	total, _ := w.Em.NumEntities()
	if total != 0 {
		t.Fatal("did not despawn all entities")
	}
	if len(w.Em.spawnSubscription.C) != 0 {
		t.Fatal("did not drain spawnSubscription channel")
	}
	if len(w.Em.despawnSubscription.C) != 0 {
		t.Fatal("did not drain despawnSubscription channel")
	}
}

func TestEntityManagerSpawnWithComponent(t *testing.T) {
	w := NewWorld(1024, 1024)

	pos := Vec2D{11, 11}
	e, _ := w.Em.Spawn(positionSpawnRequestData(pos))
	if w.Em.Components.Position[e.ID] != pos {
		t.Fatal("failed to apply component data")
	}
	if !w.Em.EntityHasComponent(e, POSITION_COMPONENT) {
		t.Fatal("failed to set entity component bit array")
	}
}

func TestEntityManagerEntitiesWithTagList(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "tag1"
	req := simpleTaggedSpawnRequestData(tag)
	w.Em.Spawn(req)
	w.Em.Update()
	tagged := w.Em.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "the chosen one"
	e, err := w.Em.UniqueTaggedEntity(tag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	req := simpleSpawnRequestData()
	_, err = w.Em.SpawnUnique(tag, req)
	if err != nil {
		t.Fatal("failed to Spawn FIRST unique entity")
	}
	e, err = w.Em.UniqueTaggedEntity(tag)
	if !(e != nil && err == nil) {
		t.Fatal("did not return unique entity")
	}
	w.Em.Update()
	_, err = w.Em.SpawnUnique(tag, req)
	if err == nil {
		t.Fatal("should not have been allowed to Spawn second unique entity")
	}
}

func TestEntityManagerTagUntagEntity(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Update()
	e := w.Em.Entities[0]
	tag := "tag1"
	w.Em.TagEntity(e, tag)
	tagged := w.Em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
	w.Em.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
}

func TestEntityManagerTagUntagEntities(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Update()
	tag := "tag1"
	w.Em.TagEntities(w.Em.GetCurrentEntities(), tag)
	tagged := w.Em.EntitiesWithTag(tag)
	if tagged.Length() != 3 {
		t.Fatal("failed to tag all entities")
	}
	w.Em.UntagEntities(w.Em.GetCurrentEntities(), tag)
	if tagged.Length() != 0 {
		t.Fatal("failed to untag all entities")
	}
}

func TestEntityManagerDeactivateActivateEntity(t *testing.T) {
	w := NewWorld(1024, 1024)

	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Update()
	e := w.Em.Entities[0]
	tag := "tag1"
	w.Em.TagEntity(e, tag)
	tagged := w.Em.EntitiesWithTag(tag)
	w.Em.Deactivate(e)
	if tagged.Length() != 0 {
		t.Fatal("entity was not removed from list after Deactivate()")
	}
	_, active := w.Em.NumEntities()
	if active != 0 {
		t.Fatal("didn't update active count")
	}
	w.Em.Activate(e)
	if tagged.Length() != 1 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
	_, active = w.Em.NumEntities()
	if active != 1 {
		t.Fatal("didn't update active count")
	}
}

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	w := NewWorld(1024, 1024)
	name := "ILoveLily"
	if w.Em.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	query := EntityQuery{
		Name: name,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			return false
		}}
	list := w.Em.GetUpdatedEntityList(query)
	if w.Em.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func TestEntityManagerString(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.Em.String()
}

func TestEntityManagerDump(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.Em.Spawn(simpleSpawnRequestData())
	w.Em.Update()
	e := w.Em.Entities[0]
	tag := "tag1"
	w.Em.TagEntity(e, tag)
	s := w.Em.Dump()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in EntityManager.Dump()")
	}
}
