/*
 * Methods needed to treat a slice of uint16's as a sorted slice of uint16's
 * using code derived from the go runtime's implementation of binary search for
 * sort.Search, specialized for uint16's
 *
 * This is used by the collision system, as provided to it by
 * EntityManager.GetSortedUpdatedEntityList, which returns a
 * SortedUpdatedEntityList, so that the collision entity ID's are always
 * in sorted order, which is needed for the triangle-packed array used there
 */

package engine

// Commentary: this is a really cute way of doing a binary search
// Returns the index to insert x at (could be len(a) if it would be new max)
func SortedUint16SliceSearch(s []uint16, x uint16) int {
	n := len(s)
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i := 0
	j := n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if s[h] < x { // i <3 u
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

func SortedUint16SliceInsert(s *[]uint16, x uint16) {
	i := SortedUint16SliceSearch(*s, x)
	*s = append(*s, 0)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = x
}

func SortedUint16SliceRemove(s *[]uint16, x uint16) {
	i := SortedUint16SliceSearch(*s, x)
	*s = append((*s)[:i], (*s)[i+1:]...)
}
