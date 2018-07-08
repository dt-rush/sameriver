package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
	"time"
)

func TestSpawn(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	em.Spawn(simpleSpawnRequestData())
	time.Sleep(16 * time.Millisecond)
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
	time.Sleep(16 * time.Millisecond)
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
	time.Sleep(16 * time.Millisecond)
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
	time.Sleep(16 * time.Millisecond)
	em.Update()
	tagged := em.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find spawned entity in EntitiesWithTag")
	}
}
