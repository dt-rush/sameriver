package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type ResettableRateLimiter struct {
	once         sync.Once
	mutex        sync.RWMutex
	delay        time.Duration
	resetCounter uint32
}

func (r *ResettableRateLimiter) Do(f func()) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	r.once.Do(func() {
		f()
		go func() {
			resetCounterBeforeSleep := atomic.LoadUint32(&r.resetCounter)
			time.Sleep(r.delay)
			if atomic.LoadUint32(&r.resetCounter) == resetCounterBeforeSleep {
				r.Reset()
			}
		}()
	})
}

func (r *ResettableRateLimiter) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	atomic.AddUint32(&r.resetCounter, 1)
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
		l.entities = append(l.entities, id)
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

func listModifier(l *EntityList) {

	boolgen := NewBoolgen()

	for {
		if boolgen.Bool() {
			l.Add()
		} else {
			l.Remove()
		}
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}
}

// Defines a kind of trianglular 2D array which allows you to store a
// ResettableRateLimiter at the intersection of each entity ID and each other
// entity ID, assuming they are indexed [i][j] where i < j
//
// For example, the table would look like the following if we had
// MAX_ENTITIES = 5, where r is a rate limiter
//
//         j
//
//     0 1 2 3 4
//    0  r r r r
//    1    r r r
// i  2      r r
//    3        r
//    4
//
type CollisionRateLimiterArray struct {
	backingArray []ResettableRateLimiter
	Arr          [][]ResettableRateLimiter
}

// Construct a new CollisionRateLimiterArray
func NewCollisionRateLimiterArray() CollisionRateLimiterArray {
	a := CollisionRateLimiterArray{
		backingArray: make([]ResettableRateLimiter,
			MAX_ENTITIES*(MAX_ENTITIES+1)/2),
		Arr: make([][]ResettableRateLimiter,
			MAX_ENTITIES),
	}
	offset := 0
	for i := 0; i < MAX_ENTITIES; i++ {
		sliceSize := MAX_ENTITIES - i
		a.Arr[i] = a.backingArray[offset : offset+sliceSize]
		offset += sliceSize
	}
	return a
}

// resets all rate limiters for the given entity
func (c *CollisionRateLimiterArray) ResetAll(id uint16) {
}

func main() {
	rand.Seed(time.Now().UnixNano())

	t0 := time.Now().UnixNano()

	// make entities
	e := NewEntityList(MAX_ENTITIES)
	for i := 0; i < rand.Intn(MAX_ENTITIES); i++ {
		e.Add()
	}
	// spawn a thread to randomly add and remove entities from the entitylist
	go listModifier(&e)

	backingArray := [*(N + 1) / 2]ResettableRateLimiter{}

	// TODO: start a loop every 16 ms that chcks all entity collisions
	// TODO: keep entities watched array sorted, so in-order traversal is
	// guaranteed. We can shift the lower value <<15 and or it to get a big
	// number which will be the key to m
}
