package engine

import (
	"testing"
)

func TestSteeringSystem(t *testing.T) {
	w := NewWorld(1024, 1024)
	ss := NewSteeringSystem()
	w.AddSystems(ss)
	e, err := w.Em.Spawn(steeringSpawnRequestData())
	vel := w.Em.Components.Velocity[e.ID]
	if err != nil {
		t.Fatal(err)
	}
	w.Update(FRAME_SLEEP_MS/2)
	if w.Em.Components.Velocity[e.ID] == vel {
		t.Fatal("failed to steer velocity")
	}
}
