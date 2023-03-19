package sameriver

import (
	"time"

	"go.uber.org/atomic"
)

// Defines a kind of trianglular 2D array which allows you to store a
// RateLimiter at the intersection of each entity ID and each other
// entity ID, assuming they are indexed [i][j] where i < j
// ("collision-indexing")
//
// For example, the table would look like the following if we had
// MAX_ENTITIES = 5, where r is a rate limiter (atomic.Uint32):
//
/*
         j

    0 1 2 3 4
    0 r r r r
    1   r r r
 i  2     r r
    3       r
    4

*/
// The rows are actually slices of a contiguous backing array, to make sure
// the rate limiters are all loaded in the same cache line
//
//             (the motiviation
//             entirely behind maintaining this weird triangle despite having
//                 a tricky algoirthm
//             to expand if we increase max entities at runtime; it would be far simpler
//             to just use a 2d array, use only the handshake cells, and it is trivial to
//             expand. Just expand empty rows and columns to the right and down.
//             The cache lines. It's all about the cache lines morty. I'd killa  man for
//             cache lines do you understand morty im a killer.
//
// The indexes to the [][] slice corresponding to the [] slice can be
// visualized like this:
/*
   i | 0 0 0 0 1 1 1 2 2 3
   j | 1 2 3 4 2 3 4 3 4 4

       r r r r r r r r r r
*/
type CollisionRateLimiterArray struct {
	capacity   int
	delay      time.Duration
	backingArr []atomic.Uint32
	arr        [][]atomic.Uint32
}

// Construct a new CollisionRateLimiterArray
func NewCollisionRateLimiterArray(nEntities int, delay time.Duration) CollisionRateLimiterArray {
	// build the backing array
	a := CollisionRateLimiterArray{
		capacity: nEntities,
		delay:    delay,
		backingArr: make(
			[]atomic.Uint32,
			nEntities*(nEntities+1)/2),
		arr: make(
			[][]atomic.Uint32,
			nEntities-1),
	}
	// now build the Arr slices which reference positions in the backing array
	//
	// consider: 3 entities
	//
	// backing arr
	// 1 2 2
	//
	// triangular slices
	// 0       1 2     [0:2]
	// 1         2     [2:3]
	offset := 0
	for i := 0; i < nEntities-1; i++ {
		sliceSize := nEntities - i - 1
		a.arr[i] = a.backingArr[offset : offset+sliceSize]
		offset += sliceSize
	}
	// return the object we've built
	return a
}

// Get the rate limiter for an i, j pair
func (a *CollisionRateLimiterArray) GetRateLimiter(i int, j int) *atomic.Uint32 {
	// for the i'th array,
	//
	//     j    |    index
	// 	 i + 1  |      0
	//   i + 2  |      1
	//   i + 3  |      2
	//    ...         ...
	//
	// so in order to map j = i + 1 to 0, j = i + 2 to 1, etc.,
	// we use j - (i+1) as the index
	return &a.arr[i][j-(i+1)]
}

// call a function in a rate-limited way
func (a *CollisionRateLimiterArray) Do(i, j int, f func()) {
	r := a.GetRateLimiter(i, j)
	if r.CompareAndSwap(0, 1) {
		f()
		go func() {
			time.Sleep(a.delay)
			r.CompareAndSwap(1, 0)
		}()
	}
}

// Reset all the rate limiters corresponding to an ID in the array (the
// entity there has been despawned)
func (a *CollisionRateLimiterArray) Reset(e *Entity) {
	// clear all where i = id
	for _, r := range a.arr[e.ID] {
		r.CompareAndSwap(1, 0)
	}
	// clear all where j = id
	for i := 0; i < e.ID; i++ {
		r := a.GetRateLimiter(i, e.ID)
		r.CompareAndSwap(1, 0)
	}
}

// this is a bit of a monstrosity.
// we need to be able to expand the capacity of the number of entities
// in total, which means expanding the backing array, moving all atomic.Uint32's
// to their new locations (sparser), and re-defining the arr subslice map
func (a *CollisionRateLimiterArray) Expand(n int) {

	Logger.Println("expanding collision rate limiter array")

	newCapacity := a.capacity + n

	newBackingArr := make([]atomic.Uint32, newCapacity*(newCapacity+1)/2)
	newArr := make([][]atomic.Uint32, newCapacity-1)

	// reposition the existing rate limiters in the new array
	/*

		old: 4
		      r r r   r r   r
		     [1 2 3] [2 3] [3]
		0    1 2 3
		1      2 3
		2        3
		new: 5
		      r r r     r r     r
		     [1 2 3 4] [2 3 4] [3 4] [4]
		0    1 2 3 4
		1      2 3 4
		2        3 4
		3          4

	*/
	offset := 0
	for i := 0; i < a.capacity-1; i++ {
		oldSliceSize := a.capacity - i - 1
		newSliceSize := newCapacity - i - 1
		for j := 0; j < oldSliceSize; j++ {
			r := a.GetRateLimiter(i, i+j+1)
			newBackingArr[offset+j] = *r
		}
		offset += newSliceSize
	}
	// update the subslice references
	offset = 0
	for i := 0; i < newCapacity-1; i++ {
		newSliceSize := newCapacity - i - 1
		newArr[i] = newBackingArr[offset : offset+newSliceSize]
		offset += newSliceSize
	}

	a.backingArr = newBackingArr
	a.arr = newArr
	a.capacity = newCapacity
}
