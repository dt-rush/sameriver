/**
  * 
  * 
  * 
  * 
**/



package utils

import (
	"time"
)

// profiler for *blocking* functions


// not very good rn

type FuncProfiler struct {
	// NO STATE
	// it's a communist utopia
	// NO MONEY
	// NO BOURGEOISIE
	
}

func (fp *FuncProfiler) Init (capacity int) {
	// NOTHING as yet
	// TODO: make this capable of averaging a function
	// that gets called repeatedly (how to catch stop?)
}

func (fp *FuncProfiler) Time (f func ()) float64 {
	t0 := time.Now().UnixNano()
	// NOTA BENE: f() must block to be timed accurately at all!
	// we're not tracking go routines here at all
	f()
	t1 := time.Now().UnixNano()
	milliseconds := float64 (t1 - t0) / 1e6
	return milliseconds
}
