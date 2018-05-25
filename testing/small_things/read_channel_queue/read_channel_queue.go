package main

import (
	"fmt"
	"math/rand"
	"time"
)

var CAPACITY = 1024
var N_WRITERS = 1024
var WRITER_SLEEP_MS = 20 * time.Millisecond

func writer(c chan int) {
	for {
		t0 := time.Now()
		c <- rand.Intn(128)
		wait := time.Since(t0).Nanoseconds() / 1e6
		if wait > 1 {
			// fmt.Printf("%d\n", wait)
		}
		// time.Sleep(WRITER_SLEEP_MS)
	}
}

func main() {
	c := make(chan int, CAPACITY)
	for i := 0; i < N_WRITERS; i++ {
		go writer(c)
	}
	for {
		bufsize := len(c)
		t0 := time.Now()
		for i := 0; i < bufsize; i++ {
			val := <-c
			val++
		}
		fmt.Printf("%d\n", time.Since(t0).Nanoseconds()/1e6)
		time.Sleep(16 * time.Millisecond)
	}
}
