package sameriver

import (
	"testing"
	"time"

	"encoding/json"
)

func testingSetupCollision() (*World, *CollisionSystem, *EventChannel, *Entity) {
	w := testingWorld()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(
		NewSpatialHashSystem(1, 1),
		cs,
	)
	// spawn the colliding entities
	testingSpawnCollision(w)
	e := testingSpawnCollision(w)
	// subscribe to collision events
	ec := w.Events.Subscribe(SimpleEventFilter("collision"))
	return w, cs, ec, e
}

func TestCollisionSystem(t *testing.T) {
	w, _, ec, e := testingSetupCollision()
	w.Update(FRAME_DURATION_INT / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	select {
	case <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 1 frame")
	}
	// move the enitity so it no longer collides
	*e.GetVec2D("Position") = Vec2D{100, 100}
	w.Update(FRAME_DURATION_INT / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	if len(ec.C) != 1 {
		t.Fatal("collision event occurred but entities were not overlapping")
	}
}

func TestCollisionSystemMany(t *testing.T) {
	w, _, ec, _ := testingSetupCollision()
	for i := 0; i < 100; i++ {
		testingSpawnCollisionRandom(w)
	}
	Logger.Printf("%d entities.", len(w.GetCurrentEntitiesSet()))
	w.SetSystemSchedule("CollisionSystem", 5)
	w.Update(FRAME_DURATION_INT / 2)
	time.Sleep(5 * FRAME_DURATION)
	w.Update(FRAME_DURATION_INT / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	select {
	case <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 1 frame")
	}
}

func TestCollisionRateLimit(t *testing.T) {
	w, cs, ec, e := testingSetupCollision()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	time.Sleep(FRAME_DURATION)
	if len(ec.C) == 4 {
		t.Fatal("collision rate-limiter didn't prevent collision duplication")
	}
	for len(ec.C) > 0 {
		<-ec.C
	}
	// wait for rate limit to die
	time.Sleep(FRAME_DURATION)
	// check if we can reset the rate lmiiter
	w.Update(FRAME_DURATION_INT / 2)
	cs.rateLimiterArray.Reset(e)
	w.Update(FRAME_DURATION_INT / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	if len(ec.C) != 4 {
		t.Fatal("collision rate-limiter reset did not allow second collision")
	}
}

func TestCollisionFilter(t *testing.T) {
	w, _, _, e := testingSetupCollision()
	coin := w.Spawn(map[string]any{
		"tags": []string{"coin"},
		"components": map[string]any{
			"Vec2D,Position": *e.GetVec2D("Position"),
			"Vec2D,Box":      *e.GetVec2D("Box"),
		}})
	predicate := func(ev Event) bool {
		c := ev.Data.(CollisionData)
		return c.This == e && c.Other == coin
	}
	ec := w.Events.Subscribe(PredicateEventFilter("collision", predicate))
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	select {
	case ev := <-ec.C:
		if !predicate(ev) {
			t.Fatal("CollisionEventFilter did not select correct event")
		}
	default:
		t.Fatal("CollisionEventFilter didn't pick up collision")
	}
}

func TestCollisionRateLimiterArrayExpand(t *testing.T) {
	a := NewCollisionRateLimiterArray(4, 500*time.Millisecond)
	a.arr[0][0].CompareAndSwap(0, 1)
	a.arr[0][1].CompareAndSwap(0, 1)
	a.arr[0][2].CompareAndSwap(0, 1)
	a.arr[1][0].CompareAndSwap(0, 1)
	a.arr[1][1].CompareAndSwap(0, 1)
	a.arr[2][0].CompareAndSwap(0, 1)

	Logger.Println("before:")
	s, _ := json.MarshalIndent(a.arr, "", "\t")
	Logger.Println(string(s))

	a.Expand(2)

	Logger.Println("after:")
	s, _ = json.MarshalIndent(a.arr, "", "\t")
	Logger.Println(string(s))

	// i, j, val
	expect := [][3]int{
		[3]int{0, 0, 1},
		[3]int{0, 1, 1},
		[3]int{0, 2, 1},
		[3]int{0, 3, 0},
		[3]int{0, 4, 0},

		[3]int{1, 0, 1},
		[3]int{1, 1, 1},
		[3]int{1, 2, 0},
		[3]int{1, 3, 0},

		[3]int{2, 0, 1},
		[3]int{2, 1, 0},
		[3]int{2, 2, 0},

		[3]int{3, 0, 0},
		[3]int{3, 1, 0},

		[3]int{4, 0, 0},
	}
	for _, e := range expect {
		if a.arr[e[0]][e[1]].Load() != uint32(e[2]) {
			t.Fatal("didn't find value in the right place after expand")
		}
	}

}