package sameriver

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestRuntimeLimitSharerLoad(t *testing.T) {
	share := NewRuntimeLimitSharer()
	share.RegisterRunner("other")
	share.RegisterRunner("loadtest")
	r := share.runnerMap["loadtest"]

	// time.Sleep doesn't like amounts < 1ms, so we scale up the time axis
	// to allow proper sleeping
	allowance_ms := 10.0
	N_EPSILON := 10
	epsilon_factor := 0.1
	N_HEAVY := 5
	heavy_factor := 0.7

	totalLoad := float64(N_EPSILON)*epsilon_factor + float64(N_HEAVY)*heavy_factor

	Logger.Printf("allowance_ms: %f", allowance_ms)
	Logger.Printf("N_EPSILON: %v", N_EPSILON)
	Logger.Printf("epsilon_factor: %v", epsilon_factor)
	Logger.Printf("N_HEAVY: %v", N_HEAVY)
	Logger.Printf("heavy_factor: %v", heavy_factor)
	Logger.Printf("total load: %f", totalLoad)

	frame := -1
	seq := make([][]string, 0)
	markRan := func(name string) {
		seq[frame] = append(seq[frame], name)
	}
	pushFrame := func() {
		frame++
		seq = append(seq, make([]string, 0))
		Logger.Printf("------------------ frame %d ----------------------", frame)
	}
	printFrame := func() {
		b, _ := json.MarshalIndent(seq[frame], "", "\t")
		Logger.Printf(string(b))
		for _, l := range r.logicUnits {
			Logger.Printf("%s: h%d", l.name, l.hotness)
		}
	}

	for i := 0; i < N_EPSILON; i++ {
		name := fmt.Sprintf("epsilon-%d", i)
		share.addLogicImmediately("loadtest",
			&LogicUnit{
				name:    name,
				worldID: i,
				f: func(dt_ms float64) {
					Logger.Printf("epsilon func sleeping %f ms", float64(time.Duration(epsilon_factor*allowance_ms*1e6)*time.Nanosecond)/1e6)
					t0 := time.Now()
					time.Sleep(time.Duration(epsilon_factor*allowance_ms*1e6) * time.Nanosecond)
					Logger.Printf("Sleep took %f ms", float64(time.Since(t0).Nanoseconds())/1e6)
					markRan(name)
				},
				active:      true,
				runSchedule: nil})
	}

	x := 0
	for i := 0; i < N_HEAVY; i++ {
		name := fmt.Sprintf("heavy-%d", i)
		share.addLogicImmediately("loadtest",
			&LogicUnit{
				name:    name,
				worldID: N_EPSILON + 1 + i,
				f: func(dt_ms float64) {
					x += 1
					markRan(name)
					time.Sleep(time.Duration(heavy_factor*allowance_ms) * time.Millisecond)
				},
				active:      true,
				runSchedule: nil})
	}

	runFrame := func(allowanceScale float64) {
		t0 := time.Now()
		if allowanceScale != 1 {
			Logger.Printf("<CONSTRICTED FRAME>")
		}
		pushFrame()
		share.Share(allowanceScale * allowance_ms)
		elapsed := float64(time.Since(t0).Nanoseconds()) / 1.0e6
		printFrame()
		Logger.Printf("            elapsed: %f ms", elapsed)
		if allowanceScale != 1 {
			Logger.Printf("</CONSTRICTED FRAME>")
		}
	}

	// since it's never run before, running the logic will set its estimate
	runFrame(1.0)

	/*
		TODO: we need better math here to account for N_EPSILON and N_HEAVY

		heavyFirstFrame := int(math.Ceil((1.0 - (float64(N_EPSILON) * 0.1)) / heavy_factor))
		Logger.Printf("Expecting %d heavies to have run in first frame", heavyFirstFrame)
		if x != heavyFirstFrame {
			t.Fatalf("Should've run %d heavies on first frame", heavyFirstFrame)
		}

		// now we try to run it again, but give it no time to run (exceeds estimate)
		runFrame(0.1)

		if x != heavyFirstFrame+1 {
			t.Fatal("should've ran one more heavy in a constricted frame")
		}
	*/

	// run a bunch of constricted frames
	for i := 0; i < 12; i++ {
		runFrame(0.333)
	}

	// run a bunch of frames
	for i := 0; i < 12; i++ {
		runFrame(1.0)
	}
}

// TODO: this doesn't actually work properly, since each invocation of a physics
// function would need ot receive the dt_ms relative to the last run. we really need for
//   - the runtimelimiter to allow looping of logics - when they come up in
//     iteration, we try to run them n times instead of just once
func TestRuntimeLimitSharerScalePhysics(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	extraPhysicsRunner := w.RuntimeSharer.RegisterRunner("extra-physics")
	physics_scale := 3
	physics_extra_calls := 0
	for i := 0; i < physics_scale; i++ {
		extraPhysics := &LogicUnit{
			name:    fmt.Sprintf("physics-extra-%d", i),
			worldID: w.IdGen.Next(),
			f: func(dt_ms float64) {
				physics_extra_calls++
				ps.Update(dt_ms)
			},
			active:      true,
			runSchedule: nil,
		}
		extraPhysicsRunner.Add(extraPhysics)
	}
	e := testingSpawnPhysics(w)
	*e.GetVec2D("Velocity") = Vec2D{1, 1}
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_DURATION_INT / 2)
	Logger.Printf("In World.Update() 0, physics ran extra %d times", physics_extra_calls)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_DURATION_INT / 2)
}
