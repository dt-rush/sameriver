package engine

// Defines a kind of trianglular 2D array which allows you to store a
// uint32 at the intersection of each entity ID and each other
// entity ID, assuming they are indexed [i][j] where i < j
// ("collision-indexing")
//
// For example, the table would look like the following if we had
// MAX_ENTITIES = 5, where u is a uint32
//
//         j
//
//     0 1 2 3 4
//    0  u u u u
//    1    u u u
// i  2      u u
//    3        u
//    4
//
// The rows are actually slices of a contiguous backing array, to make sure
// the rate limiters are nicely packed into a cache line
//
// The indexes to the [][] slice corresponding to the [] slice can be
// visualized like this:
//
// i | 0 0 0 0 1 1 1 2 2 3
// j | 1 2 3 4 2 3 4 3 4 4
//     u u u u u u u u u u
//
type CollisionUint32Array struct {
	backingArray []uint32
	Arr          [][]uint32
}

// Construct a new CollisionRateLimiterArray
func NewCollisionUint32Array() CollisionUint32Array {
	// build the backing array
	a := CollisionUint32Array{
		backingArray: make([]uint32,
			MAX_ENTITIES*(MAX_ENTITIES+1)/2),
		Arr: make([][]uint32,
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
		a.backingArray[i] = uint32(0)
	}
	// return the object we've built
	return a
}
