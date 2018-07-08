package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
)

func TestSpawn(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	em.Spawn(simpleSpawnRequestData())
	em.Update()
	if em.NumEntities() == 0 {
		t.Fatal("failed to spawn simple spawn request entity")
	}
	e := em.Entities()[0]
	if !e.Active {
		t.Fatal("spawned entity was not active")
	}
}

func TestEntityQuery(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	req := simpleSpawnRequestData()
	pos := req.Components.Position
	em.Spawn(req)
	em.Update()
	e := em.Entities()[0]
	q := engine.EntityQuery{
		"positionQuery",
		func(e *engine.EntityToken, em *engine.EntityManager) bool {
			return em.ComponentsData.Position[e.ID] == *pos
		},
	}
	if !q.Test(e, em) {
		t.Fatal("query did not return true")
	}
}

func TestEntityQueryFromTag(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	req := simpleTaggedSpawnRequestData()
	tag := req.Components.TagList.Tags[0]
	em.Spawn(req)
	em.Update()
	e := em.Entities()[0]
	q := engine.EntityQueryFromTag(tag)
	if !q.Test(e, em) {
		t.Fatal("query did not return true")
	}
}

func TestEntitiesWithTagList(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	req := simpleTaggedSpawnRequestData()
	tag := req.Components.TagList.Tags[0]
	em.Spawn(req)
	em.Update()
	tagged := em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find spawned entity in EntitiesWithTag")
	}
}

func TestEntitySpawnUnique(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	req := simpleTaggedSpawnRequestData()
	_, err := em.SpawnUnique("the chosen one", req)
	if err != nil {
		t.Fatal("failed to spawn FIRST unique entity")
	}
	em.Update()
	_, err = em.SpawnUnique("the chosen one", req)
	if err == nil {
		t.Fatal("should not have been allowed to spawn second unique entity")
	}
}

func TestTagUntagEntity(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	em.Spawn(simpleSpawnRequestData())
	em.Update()
	e := em.Entities()[0]
	tag := "tag1"
	em.TagEntity(e, tag)
	tagged := em.EntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find spawned entity in EntitiesWithTag")
	}
	em.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
}

func TestDeactivateActivateEntity(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	em.Spawn(simpleSpawnRequestData())
	em.Update()
	e := em.Entities()[0]
	tag := "tag1"
	em.TagEntity(e, tag)
}
