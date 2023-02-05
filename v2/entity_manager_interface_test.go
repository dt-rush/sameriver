package sameriver

import (
	"regexp"
	"testing"
	"time"
)

func EntityManagerInterfaceTestSpawn(
	em EntityManagerInterface, t *testing.T) {

	pos := Vec2D{11, 11}
	e, err := testingSpawnPosition(em, pos)
	if err != nil {
		t.Fatal("error on simple spawn")
	}
	if *e.GetVec2D("Position") != pos {
		t.Fatal("failed to apply component data")
	}
	total, _ := em.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	if !(e.Despawned == false &&
		e.Active == true &&
		e.World != nil) {
		t.Fatal("entity struct not populated properly")
	}
	if len(em.GetCurrentEntitiesSet()) == 0 {
		t.Fatal("entity not added to current entities set on spawn")
	}
}

func EntityManagerInterfaceTestSpawnFail(
	em EntityManagerInterface, t *testing.T) {

	for i := 0; i < MAX_ENTITIES; i++ {
		testingSpawnSimple(em)
	}
	_, err := testingSpawnSimple(em)
	if err == nil {
		t.Fatal("should have thrown error on spawnrequest when entity table " +
			"full")
	}
}

func EntityManagerInterfaceTestQueueSpawn(
	em EntityManagerInterface, t *testing.T) {

	testingQueueSpawnSimple(em)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	em.Update(FRAME_DURATION_INT / 2)
	total, _ := em.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn " +
			"request channel")
	}
}

func EntityManagerInterfaceTestQueueSpawnFull(
	em EntityManagerInterface, t *testing.T) {
	// fill up the spawnSubscription channel
	for i := 0; i < EVENT_SUBSCRIBER_CHANNEL_CAPACITY; i++ {
		testingQueueSpawnSimple(em)
	}
	// spawn two more entities (one simple, one unique)
	testingQueueSpawnSimple(em)
	testingQueueSpawnUnique(em)
	// sleep long enough for the events to appear on the channel
	time.Sleep(FRAME_DURATION)
	// update *twice*, allowing the extra events to process despite having seen
	// a full spawn subscription channel the first time
	em.Update(FRAME_DURATION_INT / 2)
	em.Update(FRAME_DURATION_INT / 2)
	total, _ := em.NumEntities()
	if total != EVENT_SUBSCRIBER_CHANNEL_CAPACITY+2 {
		t.Fatal("should have spawned entities after processing spawn " +
			"request channel")
	}
}

func EntityManagerInterfaceTestDespawn(
	em EntityManagerInterface, t *testing.T) {

	e, _ := testingSpawnSimple(em)
	em.Despawn(e)
	if e.Active {
		t.Fatal("Despawn did not deactivate entity")
	}
	if !e.Despawned {
		t.Fatal("Despawn did not set DeSpawned flag")
	}
}

func EntityManagerInterfaceTestQueueDespawn(
	em EntityManagerInterface, t *testing.T) {

	e, _ := testingSpawnSimple(em)
	em.QueueDespawn(e)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	em.Update(FRAME_DURATION_INT / 2)
	total, _ := em.NumEntities()
	if total != 0 {
		t.Fatal("should have despawned an entity after processing despawn " +
			"subscription")
	}
}

func EntityManagerInterfaceTestDespawnAll(
	em EntityManagerInterface, t *testing.T) {

	entities := make([]*Entity, 0)
	for i := 0; i < 64; i++ {
		e, _ := testingSpawnSimple(em)
		entities = append(entities, e)
	}
	for i := 0; i < 64; i++ {
		testingQueueSpawnSimple(em)
	}
	em.DespawnAll()
	total, _ := em.NumEntities()
	if total != 0 {
		t.Fatal("did not despawn all entities")
	}
	em.Update(FRAME_DURATION_INT / 2)
	total, _ = em.NumEntities()
	if total != 0 {
		t.Fatal("DespawnAll() did not discard pending spawn requests")
	}
}

func EntityManagerInterfaceTestEntityHasComponent(
	em EntityManagerInterface, t *testing.T) {

	pos := Vec2D{11, 11}
	e, _ := testingSpawnPosition(em, pos)
	if !em.EntityHasComponent(e, "Position") {
		t.Fatal("failed to set or get entity component bit array")
	}
}

func EntityManagerInterfaceTestEntitiesWithTag(
	em EntityManagerInterface, t *testing.T) {

	tag := "tag1"
	testingSpawnTagged(em, tag)
	tagged := em.UpdatedEntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func EntityManagerInterfaceTestSpawnUnique(
	em EntityManagerInterface, t *testing.T) {

	uniqueTag := "the chosen one"
	e, err := em.UniqueTaggedEntity(uniqueTag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	e, err = em.SpawnUnique(uniqueTag, []string{}, ComponentSet{})
	if err != nil {
		t.Fatal("failed to Spawn unique entity")
	}
	eRetrieved, err := em.UniqueTaggedEntity(uniqueTag)
	if !(eRetrieved == e && err == nil) {
		t.Fatal("did not return unique entity")
	}
	_, err = em.SpawnUnique(uniqueTag, []string{}, ComponentSet{})
	if err == nil {
		t.Fatal("should not have been allowed to Spawn second unique entity")
	}
}

func EntityManagerInterfaceTestTagUntagEntity(
	em EntityManagerInterface, t *testing.T) {

	e, _ := testingSpawnSimple(em)
	tag := "tag1"
	em.TagEntity(e, tag)
	tagged := em.UpdatedEntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
	if !em.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as untagged")
	}
	em.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
	if em.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as still having removed tag")
	}
}

func EntityManagerInterfaceTestTagEntities(
	em EntityManagerInterface, t *testing.T) {

	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e, _ := testingSpawnSimple(em)
		entities = append(entities, e)
	}
	em.TagEntities(entities, tag)
	for _, e := range entities {
		if !e.GetTagList("GenericTags").Has(tag) {
			t.Fatal("entity's taglist was not modified by TagEntities")
		}
	}
}

func EntityManagerInterfaceTestUntagEntities(
	em EntityManagerInterface, t *testing.T) {

	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e, _ := testingSpawnTagged(em, tag)
		entities = append(entities, e)
	}
	em.UntagEntities(entities, tag)
	for _, e := range entities {
		if e.GetTagList("GenericTags").Has(tag) {
			t.Fatal("entity's taglist was not modified by UntagEntities")
		}
	}
}

func EntityManagerInterfaceTestDeactivateActivateEntity(
	em EntityManagerInterface, t *testing.T) {

	e, _ := testingSpawnSimple(em)
	tag := "tag1"
	em.TagEntity(e, tag)
	tagged := em.UpdatedEntitiesWithTag(tag)
	em.Deactivate(e)
	if tagged.Length() != 0 {
		t.Fatal("entity was not removed from list after Deactivate()")
	}
	_, active := em.NumEntities()
	if active != 0 {
		t.Fatal("didn't update active count")
	}
	em.Activate(e)
	if tagged.Length() != 1 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
	_, active = em.NumEntities()
	if active != 1 {
		t.Fatal("didn't update active count")
	}
}

func EntityManagerInterfaceTestGetUpdatedEntityList(
	em EntityManagerInterface, t *testing.T) {

	name := "ILoveLily"
	nameToo := "ILoveLily!!!"
	list := em.GetUpdatedEntityList(
		NewEntityFilter(
			name,
			func(e *Entity) bool {
				return true
			}),
	)
	testingSpawnSimple(em)
	if list.Length() != 1 {
		t.Fatal("failed to update UpdatedEntityList")
	}
	list2 := em.GetUpdatedEntityList(
		NewEntityFilter(
			nameToo,
			func(e *Entity) bool {
				return true
			}),
	)
	if list2.Length() != 1 {
		t.Fatal("failed to created UpdatedEntityList relative to existing entities")
	}
}

func EntityManagerInterfaceTestGetSortedUpdatedEntityList(
	em EntityManagerInterface, t *testing.T) {

	list := em.GetSortedUpdatedEntityList(
		NewEntityFilter(
			"filter",
			func(e *Entity) bool {
				return true
			}),
	)
	e8 := &Entity{ID: 8, Active: true, Despawned: false}
	e0 := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e8})
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	first, _ := list.FirstEntity()
	if first.ID != 0 {
		t.Fatal("didn't insert in order")
	}
}

func EntityManagerInterfaceTestGetUpdatedEntityListByName(
	em EntityManagerInterface, t *testing.T) {

	name := "ILoveLily"
	if em.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	list := em.GetUpdatedEntityList(
		NewEntityFilter(
			name,
			func(e *Entity) bool {
				return false
			}),
	)
	if em.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func EntityManagerInterfaceTestGetCurrentEntitiesSet(
	em EntityManagerInterface, t *testing.T) {

	if !(len(em.GetCurrentEntitiesSet()) == 0) {
		t.Fatal("initially, len(GetCurrentEntitiesSet()) should be 0")
	}
	e, _ := testingSpawnSimple(em)
	if !(len(em.GetCurrentEntitiesSet()) == 1) {
		t.Fatal("after spawn, len(GetCurrentEntitiesSet()) should be 1")
	}
	em.Deactivate(e)
	if !(len(em.GetCurrentEntitiesSet()) == 1) {
		t.Fatal("after deactivate, len(GetCurrentEntitiesSet()) should be 1")
	}
	em.Despawn(e)
	if !(len(em.GetCurrentEntitiesSet()) == 0) {
		t.Fatal("after despawn, len(GetCurrentEntitiesSet()) should be 0")
	}
}

func EntityManagerInterfaceTestString(
	em EntityManagerInterface, t *testing.T) {
	if em.String() == "" {
		t.Fatal("string implementation cannot be empty string")
	}
}

func EntityManagerInterfaceTestDumpEntities(
	em EntityManagerInterface, t *testing.T) {

	e, _ := testingSpawnSimple(em)
	tag := "tag1"
	em.TagEntity(e, tag)
	s := em.DumpEntities()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in EntityManager.Dump()")
	}
}

func EntityManagerInterfaceTest(
	em EntityManagerInterface, t *testing.T) {

}
