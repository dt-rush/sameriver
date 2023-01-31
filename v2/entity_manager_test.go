package sameriver

import (
	"testing"
)

func TestEntityManagerConstruct(t *testing.T) {
	em := NewEntityManager(testingWorld())
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}

func TestEntityManagerSpawn(t *testing.T) {
	EntityManagerInterfaceTestSpawn(testingWorld().em, t)
}

func TestWorldSpawn(t *testing.T) {
	EntityManagerInterfaceTestSpawn(testingWorld(), t)
}

func TestEntityManagerSpawnFail(t *testing.T) {
	EntityManagerInterfaceTestSpawnFail(testingWorld().em, t)
}

func TestWorldSpawnFail(t *testing.T) {
	EntityManagerInterfaceTestSpawnFail(testingWorld(), t)
}

func TestEntityManagerQueueSpawn(t *testing.T) {
	EntityManagerInterfaceTestQueueSpawn(testingWorld().em, t)
}

func TestWorldQueueSpawn(t *testing.T) {
	EntityManagerInterfaceTestQueueSpawn(testingWorld(), t)
}

func TestEntityManagerDespawn(t *testing.T) {
	EntityManagerInterfaceTestDespawn(testingWorld().em, t)
}

func TestWorldDespawn(t *testing.T) {
	EntityManagerInterfaceTestDespawn(testingWorld(), t)
}

func TestEntityManagerQueueDespawn(t *testing.T) {
	EntityManagerInterfaceTestQueueDespawn(testingWorld().em, t)
}

func TestWorldQueueDespawn(t *testing.T) {
	EntityManagerInterfaceTestQueueDespawn(testingWorld(), t)
}

func TestEntityManagerDespawnAll(t *testing.T) {
	EntityManagerInterfaceTestDespawnAll(testingWorld().em, t)
}

func TestWorldDespawnAll(t *testing.T) {
	EntityManagerInterfaceTestDespawnAll(testingWorld(), t)
}

func TestEntityManagerEntityHasComponent(t *testing.T) {
	EntityManagerInterfaceTestEntityHasComponent(testingWorld().em, t)
}

func TestWorldEntityHasComponent(t *testing.T) {
	EntityManagerInterfaceTestEntityHasComponent(testingWorld(), t)
}

func TestEntityManagerEntitiesWithTag(t *testing.T) {
	EntityManagerInterfaceTestEntitiesWithTag(testingWorld().em, t)
}

func TestWorldEntitiesWithTag(t *testing.T) {
	EntityManagerInterfaceTestEntitiesWithTag(testingWorld(), t)
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	EntityManagerInterfaceTestSpawnUnique(testingWorld().em, t)
}

func TestWorldSpawnUnique(t *testing.T) {
	EntityManagerInterfaceTestSpawnUnique(testingWorld(), t)
}

func TestEntityManagerTagUntagEntity(t *testing.T) {
	EntityManagerInterfaceTestTagUntagEntity(testingWorld().em, t)
}

func TestWorldTagUntagEntity(t *testing.T) {
	EntityManagerInterfaceTestTagUntagEntity(testingWorld(), t)
}

func TestEntityManagerTagEntities(t *testing.T) {
	EntityManagerInterfaceTestTagEntities(testingWorld().em, t)
}

func TestWorldTagEntities(t *testing.T) {
	EntityManagerInterfaceTestTagEntities(testingWorld(), t)
}

func TestEntityManagerUntagEntities(t *testing.T) {
	EntityManagerInterfaceTestUntagEntities(testingWorld().em, t)
}

func TestWorldUntagEntities(t *testing.T) {
	EntityManagerInterfaceTestUntagEntities(testingWorld(), t)
}

func TestEntityManagerDeactivateActivateEntity(t *testing.T) {
	EntityManagerInterfaceTestDeactivateActivateEntity(testingWorld().em, t)
}

func TestWorldDeactivateActivateEntity(t *testing.T) {
	EntityManagerInterfaceTestDeactivateActivateEntity(testingWorld(), t)
}

func TestEntityManagerGetUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityList(testingWorld().em, t)
}

func TestWorldGetUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityList(testingWorld(), t)
}

func TestEntityManagerGetSortedUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetSortedUpdatedEntityList(testingWorld().em, t)
}

func TestWorldGetSortedUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetSortedUpdatedEntityList(testingWorld(), t)
}

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityListByName(testingWorld().em, t)
}

func TestWorldGetUpdatedEntityListByName(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityListByName(testingWorld(), t)
}

func TestEntityManagerGetCurrentEntitiesSet(t *testing.T) {
	EntityManagerInterfaceTestGetCurrentEntitiesSet(testingWorld().em, t)
}

func TestWorldGetCurrentEntitiesSet(t *testing.T) {
	EntityManagerInterfaceTestGetCurrentEntitiesSet(testingWorld(), t)
}

func TestEntityManagerString(t *testing.T) {
	EntityManagerInterfaceTestString(testingWorld().em, t)
}

func TestWorldString(t *testing.T) {
	EntityManagerInterfaceTestString(testingWorld(), t)
}

func TestEntityManagerDumpEntities(t *testing.T) {
	EntityManagerInterfaceTestDumpEntities(testingWorld().em, t)
}

func TestWorldDumpEntities(t *testing.T) {
	EntityManagerInterfaceTestDumpEntities(testingWorld(), t)
}
