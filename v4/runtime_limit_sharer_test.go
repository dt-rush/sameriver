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
	r := share.RunnerMap["loadtest"]

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
		share.RunnerMap["loadtest"].addLogicImmediately(
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
		share.RunnerMap["loadtest"].addLogicImmediately(
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

/*
output example:

allowance_ms: 1600.000000
N_EPSILON: 1048
worksleep: 10.000000
totalLoad: 6.550000
------------------ frame 0 ----------------------
no-overhead avg hotness expected: h0.153
realised avg hotness: h0.150
ratio: 0.981250
elapsed: 1602.845728 ms

note that the realised avg hotness 0.150 is not quite the theoretical
1 / totalLoad. Because totalLoad is calculated based on gapless division
by the worksleep amount. But Really, the worksleep is bracketed by overhead.

1 - 0.981250 = about 1.8% overhead, not a bad price for the exchange of
getting (attempted, modulo roundrobin necessities) runtime limiting
*/
func TestRuntimeLimitSharerCapacity(t *testing.T) {
	share := NewRuntimeLimitSharer()
	share.RegisterRunner("capacitytest")
	r := share.RunnerMap["capacitytest"]

	// time.Sleep doesn't like amounts < 1ms, so we scale up the time axis
	// to allow proper sleeping
	allowance_ms := 1600.0
	N_EPSILON := 1048
	worksleep := 10.0

	totalLoad := float64(N_EPSILON) * worksleep / allowance_ms

	Logger.Printf("allowance_ms: %f", allowance_ms)
	Logger.Printf("N_EPSILON: %d", N_EPSILON)
	Logger.Printf("worksleep: %f", worksleep)
	Logger.Printf("totalLoad: %f", totalLoad)

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
		// b, _ := json.MarshalIndent(seq[frame], "", "\t")
		// Logger.Printf(string(b))
		hotnessSum := 0.0
		for _, l := range r.logicUnits {
			hotnessSum += float64(l.hotness)
		}
		ideal := 1 / totalLoad
		realised := hotnessSum / float64(len(r.logicUnits))
		Logger.Printf("no-overhead avg hotness expected: h%.3f", ideal)
		Logger.Printf("realised avg hotness: h%.3f", realised)
		Logger.Printf("ratio: %f", realised/ideal)
	}

	// set true to observe that we worksleep longer than 1 ms
	const observeSleep = false

	for i := 0; i < N_EPSILON; i++ {
		name := fmt.Sprintf("epsilon-%d", i)
		share.RunnerMap["capacitytest"].addLogicImmediately(
			&LogicUnit{
				name:    name,
				worldID: i,
				f: func(dt_ms float64) {
					var t0 time.Time
					if observeSleep {
						t0 = time.Now()
					}
					time.Sleep(time.Duration(worksleep*1e6) * time.Nanosecond)
					if observeSleep {
						Logger.Printf("elapsed: %f ms", float64(time.Since(t0).Nanoseconds())/1e6)
					}
					markRan(name)
				},
				active:      true,
				runSchedule: nil})
	}

	runFrame := func(allowanceScale float64) {
		t0 := time.Now()
		pushFrame()
		share.Share(allowanceScale * allowance_ms)
		elapsed := float64(time.Since(t0).Nanoseconds()) / 1.0e6
		printFrame()
		Logger.Printf("            elapsed: %f ms", elapsed)
	}

	// since it's never run before, running the logic will set its estimate
	runFrame(1.0)
}
