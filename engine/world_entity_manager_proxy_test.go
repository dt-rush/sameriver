package engine

import (
	"regexp"
	"testing"
	"time"
)

func TestWorldEntityManagerProxySpawnDespawn(t *testing.T) {
	w := testingWorld()
	req := simpleSpawnRequest()
	w.Spawn(req.Tags, req.Components)
	total, _ := w.em.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	e := w.em.entities[0]
	if !e.Active {
		t.Fatal("Spawned entity was not active")
	}
	w.Despawn(e)
	if e.Active {
		t.Fatal("DeSpawn did not deactivate entity")
	}
	if !e.Despawned {
		t.Fatal("DeSpawn did not set Despawned flag")
	}
}

func TestWorldEntityManagerProxySpawnDespawnEvent(t *testing.T) {
	w := testingWorld()
	w.QueueSpawn(simpleSpawnRequest())
	time.Sleep(FRAME_SLEEP)
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ := w.em.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn subscription")
	}
	e := w.em.entities[0]
	w.QueueDespawn(DespawnRequestData{Entity: e})
	time.Sleep(FRAME_SLEEP)
	w.Update(FRAME_SLEEP_MS / 2)
	total, _ = w.em.NumEntities()
	if total != 0 {
		t.Fatal("should have despawned entity after processing despawn subscription")
	}
}

func TestWorldEntityManagerProxySpawnUnique(t *testing.T) {
	w := testingWorld()
	tag := "the chosen one"
	e, err := w.UniqueTaggedEntity(tag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	req := simpleSpawnRequest()
	_, err = w.SpawnUnique(tag, req.Tags, req.Components)
	if err != nil {
		t.Fatal("failed to Spawn FIRST unique entity")
	}
	e, err = w.UniqueTaggedEntity(tag)
	if !(e != nil && err == nil) {
		t.Fatal("did not return unique entity")
	}
	_, err = w.SpawnUnique(tag, req.Tags, req.Components)
	if err == nil {
		t.Fatal("should not have been allowed to Spawn second unique entity")
	}
}

func TestWorldEntityManagerProxyEntitiesWithTag(t *testing.T) {
	w := testingWorld()
	tag := "tag1"
	req := simpleTaggedSpawnRequest(tag)
	w.Spawn(req.Tags, req.Components)
	tagged := w.EntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestWorldEntityManagerProxyDeactivateActivate(t *testing.T) {
	w := testingWorld()
	req := simpleSpawnRequest()
	w.Spawn(req.Tags, req.Components)
	e := w.em.entities[0]
	tag := "tag1"
	w.TagEntity(e, tag)
	tagged := w.EntitiesWithTag(tag)
	w.Deactivate(e)
	if tagged.Length() != 0 {
		t.Fatal("entity was not removed from list after Deactivate()")
	}
	_, active := w.NumEntities()
	if active != 0 {
		t.Fatal("didn't update active count")
	}
	w.Activate(e)
	if tagged.Length() != 1 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
	_, active = w.NumEntities()
	if active != 1 {
		t.Fatal("didn't update active count")
	}
}

func TestWorldEntityManagerProxyGetUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	name := "ILoveLily"
	list := w.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return true
		}),
	)
	req := simpleSpawnRequest()
	w.Spawn(req.Tags, req.Components)
	if list.Length() != 1 {
		t.Fatal("failed to update UpdatedEntityList")
	}
	list2 := w.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return true
		}),
	)
	if list2.Length() != 1 {
		t.Fatal("failed to created UpdatedEntityList relative to existing entities")
	}
}

func TestWorldEntityManagerProxyGetSortedUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	list := w.GetSortedUpdatedEntityList(NewEntityFilter("filter",
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

func TestWorldEntityManagerProxyGetUpdatedEntityListByName(t *testing.T) {
	w := testingWorld()
	name := "ILoveLily"
	if w.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	list := w.GetUpdatedEntityList(NewEntityFilter(name,
		func(e *Entity) bool {
			return false
		}),
	)
	if w.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func TestWorldEntityManagerProxyEntityHasComponent(t *testing.T) {
	w := testingWorld()
	pos := Vec2D{11, 11}
	req := positionSpawnRequest(pos)
	e, _ := w.Spawn(req.Tags, req.Components)
	if w.em.components.Position[e.ID] != pos {
		t.Fatal("failed to apply component data")
	}
	if !w.EntityHasComponent(e, POSITION_COMPONENT) {
		t.Fatal("failed to set entity component bit array")
	}
}

func TestWorldEntityManagerProxyEntityHasTag(t *testing.T) {
	w := testingWorld()
	tag := "tag1"
	req := simpleTaggedSpawnRequest(tag)
	e, _ := w.Spawn(req.Tags, req.Components)
	if !w.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag was false for tagged entity")
	}
}

func TestWorldEntityManagerProxyTagEntity(t *testing.T) {
	w := testingWorld()
	req := simpleSpawnRequest()
	e, _ := w.Spawn(req.Tags, req.Components)
	tag := "tag1"
	w.TagEntity(e, tag)
	if !w.em.components.TagList[e.ID].Has(tag) {
		t.Fatal("entity's taglist was not modified by TagEntity")
	}
}

func TestWorldEntityManagerProxyTagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		req := simpleSpawnRequest()
		e, _ := w.Spawn(req.Tags, req.Components)
		entities = append(entities, e)
	}
	w.TagEntities(entities, tag)
	for _, e := range entities {
		if !w.em.components.TagList[e.ID].Has(tag) {
			t.Fatal("entity's taglist was not modified by TagEntities")
		}
	}
}

func TestWorldEntityManagerProxyUntagEntity(t *testing.T) {
	w := testingWorld()
	tag := "tag1"
	req := simpleTaggedSpawnRequest(tag)
	e, _ := w.Spawn(req.Tags, req.Components)
	w.UntagEntity(e, tag)
	if w.em.components.TagList[e.ID].Has(tag) {
		t.Fatal("entity's taglist was not modified by UntagEntity")
	}
}

func TestWorldEntityManagerProxyUntagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		req := simpleTaggedSpawnRequest(tag)
		e, _ := w.Spawn(req.Tags, req.Components)
		entities = append(entities, e)
	}
	w.UntagEntities(entities, tag)
	for _, e := range entities {
		if w.em.components.TagList[e.ID].Has(tag) {
			t.Fatal("entity's taglist was not modified by UntagEntities")
		}
	}
}

func TestWorldEntityManagerProxyNumEntities(
	t *testing.T) {
	w := testingWorld()
	total, active := w.NumEntities()
	if !(total == 0 && active == 0) {
		t.Fatal("somehow total and active were nonzero")
	}
	req := simpleSpawnRequest()
	e, _ := w.Spawn(req.Tags, req.Components)
	total, active = w.NumEntities()
	if !(total == 1 && active == 1) {
		t.Fatal("total or active not updated after one spawn")
	}
	w.Deactivate(e)
	total, active = w.NumEntities()
	if !(total == 1 && active == 0) {
		t.Fatal("active count not updated or deactivate despawned")
	}
	w.Despawn(e)
	total, active = w.NumEntities()
	if !(total == 0 && active == 0) {
		t.Fatal("did not despawn properly")
	}
}

func TestWorldEntityManagerProxyGetCurrentEntities(t *testing.T) {
	w := testingWorld()
	if !(len(w.GetCurrentEntities()) == 0) {
		t.Fatal("initially, len(GetCurrentEntities()) should be 0")
	}
	req := simpleSpawnRequest()
	e, _ := w.Spawn(req.Tags, req.Components)
	if !(len(w.GetCurrentEntities()) == 1) {
		t.Fatal("after spawn, len(GetCurrentEntities()) should be 1")
	}
	w.Deactivate(e)
	if !(len(w.GetCurrentEntities()) == 1) {
		t.Fatal("after deactivate, len(GetCurrentEntities()) should be 1")
	}
	w.Despawn(e)
	if !(len(w.GetCurrentEntities()) == 0) {
		t.Fatal("after despawn, len(GetCurrentEntities()) should be 0")
	}
}

func TestWorldEntityManagerProxyEntityManagerString(t *testing.T) {
	w := testingWorld()
	w.EntityManagerString()
}

func TestWorldEntityManagerProxyDumpEntities(
	t *testing.T) {
	w := testingWorld()
	req := simpleSpawnRequest()
	w.Spawn(req.Tags, req.Components)
	e := w.em.entities[0]
	tag := "tag1"
	w.TagEntity(e, tag)
	s := w.DumpEntities()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in World.DumpEntities()")
	}
}
