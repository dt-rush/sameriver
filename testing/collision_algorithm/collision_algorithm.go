package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RateLimiter struct {
	once  sync.Once
	mutex sync.RWMutex
	delay time.Duration
}

func (r *RateLimiter) Invoke(f func()) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	f()
	time.Sleep(r.delay)
	r.once = sync.Once{}
}

type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func NewBoolgen() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

func removeUint16FromSlice(x uint16, slice *[]uint16) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

const MAX_ENTITIES = 128
const msRateLimit = 100

type EntityList struct {
	entities     []uint16
	availableIDs []uint16
	lock         sync.RWMutex
}

func NewEntityList(max_entities int) EntityList {
	availableIDs := make([]uint16, max_entities)
	for i := 0; i < max_entities; i++ {
		availableIDs[i] = uint16(i)
	}
	return EntityList{availableIDs: availableIDs}
}

func (l *EntityList) Add() {
	l.lock.Lock()
	defer l.lock.Lock()

	if len(l.availableIDs) > 0 {
		last_ix := len(l.availableIDs) - 1
		id := l.availableIDs[last_ix]
		l.availableIDs = l.availableIDs[last_ix:]
	}
}

func (l *EntityList) Remove() {
	l.lock.Lock()
	defer l.lock.Lock()

	last_ix := len(l.entities) - 1
	remove_ix := rand.Intn(len(l.entities))
	id := l.entities[remove_ix]
	l.availableIDs = append(l.availableIDs, id)
	l.entities[remove_ix] = l.entities[last_ix]
	l.entities = l.entities[last_ix:]
}

func listModifier(l *EntityList, stop chan (bool)) {

	boolgen := NewBoolgen()

	insert := func() {
		l.entities = append(l.entities, uint16(rand.Intn(MAX_ENTITIES)))
	}

	remove := func() {

	}

modifyloop:
	for {
		select {
		case <-stop:
			break modifyloop
		default:
			if boolgen.Bool() {
				l.Add()
			} else {
				l.Remove()
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m := make(map[uint32]RateLimiter)
	t0 := time.Now().UnixNano()
	fmt.Println("Allocating ID's...")
	// make entities
	e := NewEntityList()
	for i := uint16(0); i < rand.Intn(MAX_ENTITIES); i++ {
		e.insertRandom()
	}
	// spawn a thread to randomly add and remove entities from the entitylist
	go func() {

	}()
	// TODO: start a loop every 16 ms that chcks all entity collisions
	// TODO: keep entities watched array sorted, so in-order traversal is
	// guaranteed. We can shift the lower value <<15 and or it to get a big
	// number which will be the key to m
}
