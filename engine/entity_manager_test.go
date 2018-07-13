package engine

import (
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
	if w.em.NumEntities() == 0 {
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

func TestEntityManagerEntitiesWithTagList(t *testing.T) {
	w := NewWorld(1024, 1024)

	req := simpleTaggedSpawnRequestData()
	tag := req.Components.TagList.Tags[0]
	w.em.Spawn(req)
	w.em.Update()
	tagged := w.em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	w := NewWorld(1024, 1024)

	req := simpleTaggedSpawnRequestData()
	_, err := w.em.SpawnUnique("the chosen one", req)
	if err != nil {
		t.Fatal("failed to Spawn FIRST unique entity")
	}
	w.em.Update()
	_, err = w.em.SpawnUnique("the chosen one", req)
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
	w.em.Activate(e)
	if tagged.Length() == 0 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
}

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	w := NewWorld(1024, 1024)

	name := "ILoveLily"
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

func TestEntityManagerToString(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.em.String()
}
