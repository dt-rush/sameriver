package sameriver

import (
	"math/rand"
)

func testingSimpleSpawnRequest() SpawnRequestData {
	return SpawnRequestData{Tags: []string{}, Components: ComponentSet{}}
}

func testingQueueSpawnSimple(em EntityManagerInterface) {
	em.QueueSpawn([]string{}, ComponentSet{})
}

func testingQueueSpawnUnique(em EntityManagerInterface) {
	em.QueueSpawnUnique("unique", []string{}, ComponentSet{})
}

func testingSpawnSimple(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{}, ComponentSet{})
}

func testingSpawnPosition(
	em EntityManagerInterface, pos Vec2D) (*Entity, error) {
	return em.Spawn([]string{}, MakeComponentSet(map[string]interface{}{
		"Vec2D,Position": pos,
	}))
}

func testingSpawnTagged(
	em EntityManagerInterface, tag string) (*Entity, error) {
	return em.Spawn([]string{tag}, ComponentSet{})
}

func testingSpawnSpatial(
	em EntityManagerInterface, pos Vec2D, box Vec2D) (*Entity, error) {
	return em.Spawn([]string{},
		MakeComponentSet(map[string]interface{}{
			"Vec2D,Position": pos,
			"Vec2D,Box":      box,
		}))
}

func testingSpawnCollision(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{},
		MakeComponentSet(map[string]interface{}{
			"Vec2D,Position": Vec2D{10, 10},
			"Vec2D,Box":      Vec2D{4, 4},
		}))
}

func testingSpawnCollisionRandom(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{},
		MakeComponentSet(map[string]interface{}{
			"Vec2D,Position": Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			"Vec2D,Box":      Vec2D{5, 5},
		}))
}

func testingSpawnSteering(em EntityManagerInterface) (*Entity, error) {
	mass := 3.0
	maxV := 3.0
	return em.Spawn([]string{},
		MakeComponentSet(map[string]interface{}{
			"Vec2D,Position":       Vec2D{0, 0},
			"Vec2D,Velocity":       Vec2D{0, 0},
			"Float64,MaxVelocity":  maxV,
			"Vec2D,MovementTarget": Vec2D{1, 1},
			"Vec2D,Steer":          Vec2D{0, 0},
			"Float64,Mass":         mass,
		}))
}

func testingSpawnPhysics(em EntityManagerInterface) (*Entity, error) {
	mass := 3.0
	return em.Spawn([]string{},
		MakeComponentSet(map[string]interface{}{
			"Vec2D,Position": Vec2D{10, 10},
			"Vec2D,Velocity": Vec2D{0, 0},
			"Vec2D,Box":      Vec2D{1, 1},
			"Float64,Mass":   mass,
		}))
}
