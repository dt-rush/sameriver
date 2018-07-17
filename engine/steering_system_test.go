package engine

import (
	"testing"
)

func TestSteeringSystem(t *testing.T) {
	w := testingWorld()
	ss := NewSteeringSystem()
	w.AddSystems(ss)
	e, err := w.em.spawn(steeringSpawnRequest())
	vel := w.em.components.Velocity[e.ID]
	if err != nil {
		t.Fatal(err)
	}
	w.Update(FRAME_SLEEP_MS / 2)
	if w.em.components.Velocity[e.ID] == vel {
		t.Fatal("failed to steer velocity")
	}
}
