package engine

import (
	"github.com/dt-rush/sameriver/engine/utils"
	"time"
)

// Defines a kind of trianglular 2D array which allows you to store a
// ResettableRateLimiter at the intersection of each entity ID and each other
// entity ID, assuming they are indexed [i][j] where i < j
// ("collision-indexing")
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
// The rows are actually slices of a contiguous backing array, to make sure
// the rate limiters are all loaded in the same cache line
//
// The indexes to the [][] slice corresponding to the [] slice can be
// visualized like this:
//
// i | 0 0 0 0 1 1 1 2 2 3
// j | 1 2 3 4 2 3 4 3 4 4
//     r r r r r r r r r r
//
type CollisionRateLimiterArray struct {
	backingArray []*utils.ResettableRateLimiter
	Arr          [][]*utils.ResettableRateLimiter
}

// Construct a new CollisionRateLimiterArray
func NewCollisionRateLimiterArray() CollisionRateLimiterArray {
	// build the backing array
	a := CollisionRateLimiterArray{
		backingArray: make(
			[]*utils.ResettableRateLimiter,
			MAX_ENTITIES*(MAX_ENTITIES+1)/2),
		Arr: make(
			[][]*utils.ResettableRateLimiter,
			MAX_ENTITIES),
	}
	// build the Arr slices which reference positions in the backing array
	offset := 0
	for i := 0; i < MAX_ENTITIES; i++ {
		sliceSize := MAX_ENTITIES - i
		a.Arr[i] = a.backingArray[offset : offset+sliceSize]
		offset += sliceSize
	}
	// create the rate limiters
	for i := 0; i < len(a.backingArray); i++ {
		a.backingArray[i] = utils.NewResettableRateLimiter(
			COLLISION_RATELIMIT_TIMEOUT_MS * time.Millisecond)
	}
	// return the object we've built
	return a
}

// Get the rate limiter for an i, j pair
func (a *CollisionRateLimiterArray) GetRateLimiter(
	i int, j int) *utils.ResettableRateLimiter {

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
	return a.Arr[i][j-(i+1)]
}

// Reset all the rate limiters corresponding to an ID in the array (the
// entity there has been despawned)
func (a *CollisionRateLimiterArray) Reset(entity *EntityToken) {
	// clear all where i = id
	for _, r := range a.Arr[entity.ID] {
		r.Reset()
	}
	// clear all where j = id
	for i := 0; i < entity.ID; i++ {
		a.GetRateLimiter(i, entity.ID).Reset()
	}
}
