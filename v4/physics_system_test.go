package sameriver

import (
	"testing"
	"time"
)

func TestPhysicsSystemWithGranularity(t *testing.T) {
	// normal setup
	w := testingWorld()
	p := NewPhysicsSystem()
	w.RegisterSystems(p)
	e := testingSpawnPhysics(w)
	*e.GetVec2D(VELOCITY) = Vec2D{1, 1}
	pos := e.GetVec2D(POSITION)
	pos0 := *pos
	// granular setup
	wg := testingWorld()
	pg := NewPhysicsSystemWithGranularity(4)
	wg.RegisterSystems(pg)
	eg := testingSpawnPhysics(wg)
	*eg.GetVec2D(VELOCITY) = Vec2D{1, 1}
	posg := eg.GetVec2D(POSITION)
	posg0 := *posg

	// simulate constant load of other logics with a ratio
	// value of 8 means 1/8th to physics
	physicsTimeShareReciprocal := 8.0

	// run a frame of same allowance ms for both normal and granular
	runFrame := func() {
		w.Update(FRAME_MS / physicsTimeShareReciprocal)
		wg.Update(FRAME_MS / physicsTimeShareReciprocal)
		time.Sleep(FRAME_DURATION)
	}

	observePos := func() {
		Logger.Printf("normal pos: %v", *pos)
		Logger.Printf("granular pos: %v", *posg)
		Logger.Printf("position.x ratio: %f", pos.X/posg.X)
	}

	// Frame 0
	runFrame()
	Logger.Println("after Update at t=0, hotness of physics update:")
	// observe, in the below, the hotness are basically the same for
	// granularity 1 as granularity 4, insane
	for _, l := range w.RuntimeSharer.RunnerMap["systems"].logicUnits {
		Logger.Printf("normal %s: h%d", l.name, l.hotness)
	}
	for _, l := range wg.RuntimeSharer.RunnerMap["systems"].logicUnits {
		Logger.Printf("granular %s: h%d", l.name, l.hotness)
	}
	observePos()

	// Frame 1
	Logger.Println("TEST FRAME 2")
	runFrame()
	if *pos == pos0 {
		t.Fatal("failed to update position")
	}
	if *posg == posg0 {
		t.Fatal("failed to update position in granular")
	}
	// as of this comment, 2023-03-18, observe that the numeric result is different;
	// this is *at least* because
	// the physics update is getting slightly different unstable dts, due ultimately to
	// the runtimelimiter passing different dt_ms to the logicunit each time it runs
	// based on wall time since it last scheduled it, which can vary over a single
	// frame as it tries to pack in the time and repeatedly polls time since last
	// run.
	observePos()

	return
	// let's observe the behaviour over a longer term
	for i := 0; i < 30; i++ {
		runFrame()
	}
	Logger.Printf("+ 30 frames:")
	// observe that they're actually closer than after the 1st or 2nd frame
	observePos()
	for i := 0; i < 60; i++ {
		runFrame()
	}
	Logger.Printf("+ 60 frames:")
	// observe that they're actually closer than after the 1st or 2nd frame
	observePos()
}

func TestPhysicsSystemMotion(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)
	*e.GetVec2D(VELOCITY) = Vec2D{1, 1}
	pos := *e.GetVec2D(POSITION)
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_MS / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_MS / 2)
	if *e.GetVec2D(POSITION) == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemMany(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	for i := 0; i < 500; i++ {
		testingSpawnPhysics(w)
	}
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_MS / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_MS / 2)
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)
	directions := []Vec2D{
		Vec2D{100, 0},
		Vec2D{-100, 0},
		Vec2D{0, 100},
		Vec2D{0, -100},
	}
	worldCenter := Vec2D{w.Width / 2, w.Height / 2}
	worldTopRight := Vec2D{w.Width, w.Height}
	pos := e.GetVec2D(POSITION)
	box := e.GetVec2D(BOX)
	vel := e.GetVec2D(VELOCITY)
	for _, d := range directions {
		*pos = Vec2D{512, 512}
		*vel = d
		for i := 0; i < 64; i++ {
			w.Update(FRAME_MS / 2)
			time.Sleep(1 * time.Millisecond)
		}
		if !RectWithinRect(*pos, *box, worldCenter, worldTopRight) {
			t.Fatalf("traveling with velocity %v placed entity "+
				"outside world (at position %v, box %v)", *vel, *pos, *box)
		}
	}
}
