package engine

import (
	"testing"
)

func TestSteeringSystem(t *testing.T) {
	w := NewWorld(1024, 1024)
	e, err := w.em.Spawn(steeringSpawnRequestData())
	vel := w.em.Components.Velocity[e.ID]
	if err != nil {
		t.Fatal(err)
	}
	w.ss.Update(FRAME_SLEEP_MS)
	if w.em.Components.Velocity[e.ID] == vel {
		t.Fatal("failed to steer velocity")
	}
}
