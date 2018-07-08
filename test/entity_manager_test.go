package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
	"time"
)

func simpleSpawnRequestData() engine.SpawnRequestData {
	return engine.SpawnRequestData{
		Components: engine.ComponentSet{
			Position: &engine.Vec2D{0, 0},
		},
	}
}

func simpleTaggedSpawnRequestData() engine.SpawnRequestData {
	return engine.SpawnRequestData{
		Components: engine.ComponentSet{
			Position: &engine.Vec2D{0, 0},
			TagList:  &engine.TagList{Tags: []string{"tag1"}},
		},
	}
}

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

func TestTagListSpawnInserted(t *testing.T) {
	ev := engine.NewEventBus()
	em := engine.NewEntityManager(ev)
	req := simpleTaggedSpawnRequestData()
	tag := req.Components.TagList.Tags[0]
	em.Spawn(req)
	time.Sleep(16 * time.Millisecond)
	em.Update()
	if em.NumEntities() == 0 {
		t.Fatal("failed to spawn simple tagged entity")
	}
	tagged := em.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to put spawned entity in list of entities with tag")
	}
}
