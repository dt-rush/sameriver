/*
 * Methods needed to treat a slice of *Entitys as a sorted slice of
 * *Entitys using code derived from the go runtime's implementation of
 * binary search for sort.Search, specialized for *Entitys
 *
 * This is used by the collision system, as provided to it by
 * EntityManager.GetSortedUpdatedEntityList, which returns a
 * SortedUpdatedEntityList, so that the collision entity ID's are always
 * in sorted order, which is needed for the triangle-packed array used there
 */

package engine

// Returns the index to insert x at (could be len(a) if it would be new max)
func SortedEntitySliceSearch(s []*Entity, x *Entity) int {
	n := len(s)
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i := 0
	j := n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if s[h].ID < x.ID { // i <3 u
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

func SortedEntitySliceInsertIfNotPresent(
	s *[]*Entity, x *Entity) bool {

	i := SortedEntitySliceSearch(*s, x)
	// put the element at the end of the array
	*s = append(*s, nil)
	// if the insertion point was not the end,
	if i != len(*s) {
		// and the element is already there
		if (*s)[i] == x {
			// shrink the slice by 1 since we appended for nothing
			*s = (*s)[:len(*s)-1]
			return false
		} else {
			// else shift everything up and place the entity in its proper
			// position (the appended copy was overwritten by the shift)
			copy((*s)[i+1:], (*s)[i:])
			(*s)[i] = x
		}
	}
	return true
}

func SortedEntitySliceRemove(s *[]*Entity, x *Entity) {
	i := SortedEntitySliceSearch(*s, x)
	found := (i != len(*s) && (*s)[i] == x)
	if found {
		*s = append((*s)[:i], (*s)[i+1:]...)
	}
}
