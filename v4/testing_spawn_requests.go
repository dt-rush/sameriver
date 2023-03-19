package sameriver

import (
	"math/rand"
)

func testingQueueSpawnSimple(em EntityManagerInterface) {
	em.QueueSpawn(nil)
}

func testingQueueSpawnUnique(em EntityManagerInterface) {
	em.QueueSpawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnUnique(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnSimple(em EntityManagerInterface) *Entity {
	return em.Spawn(nil)
}

func testingSpawnPosition(
	em EntityManagerInterface, pos Vec2D) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]any{
			"Vec2D,Position": pos,
		}})
}

func testingSpawnTagged(
	em EntityManagerInterface, tag string) *Entity {
	return em.Spawn(map[string]any{
		"tags": []string{tag},
	})
}

func testingSpawnSpatial(
	em EntityManagerInterface, pos Vec2D, box Vec2D) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]any{
			"Vec2D,Position": pos,
			"Vec2D,Box":      box,
		}})
}

func testingSpawnCollision(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]interface{}{
			"Vec2D,Position": Vec2D{10, 10},
			"Vec2D,Box":      Vec2D{4, 4},
		}})
}

func testingSpawnCollisionRandom(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]interface{}{
			"Vec2D,Position": Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			"Vec2D,Box":      Vec2D{5, 5},
		}})
}

func testingSpawnSteering(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]interface{}{
			"Vec2D,Position":       Vec2D{0, 0},
			"Vec2D,Velocity":       Vec2D{0, 0},
			"Vec2D,Acceleration":   Vec2D{0, 0},
			"Float64,MaxVelocity":  3.0,
			"Vec2D,MovementTarget": Vec2D{1, 1},
			"Vec2D,Steer":          Vec2D{0, 0},
			"Float64,Mass":         3.0,
		}})
}

func testingSpawnPhysics(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]interface{}{
			"Vec2D,Position":     Vec2D{10, 10},
			"Vec2D,Velocity":     Vec2D{0, 0},
			"Vec2D,Acceleration": Vec2D{0, 0},
			"Vec2D,Box":          Vec2D{1, 1},
			"Float64,Mass":       3.0,
		}})
}
