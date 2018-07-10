package engine

import (
	"testing"
)

func TestSteeringSystem(t *testing.T) {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	ms := NewSteeringSystem(em)
	e, err := em.spawn(movementSpawnRequestData())
	vel := em.Components.Velocity[e.ID]
	if err != nil {
		t.Fatal(err)
	}
	ms.Update(FRAME_SLEEP_MS)
	if em.Components.Velocity[e.ID] == vel {
		t.Fatal("failed to steer velocity")
	}
}
