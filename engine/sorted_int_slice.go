package engine

// Returns the index to insert x at (could be len(a) if it would be new max)
func SortedIntSliceSearch(s []int, x int) int {
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

func SortedIntSliceInsertIfNotPresent(s *[]int, x int) bool {
	i := SortedIntSliceSearch(*s, x)
	// put the element at the end of the array
	// if the insertion point was not the end,
	if i != len(*s) {
		// and the element is already there
		if (*s)[i] == x {
			// shrink the slice by 1 since we appended for nothing
			*s = (*s)[:len(*s)-1]
			return false
		} else {
			// else shift everything up and place the entity in its proper
			// position (the 0 gets overwritten by the shift)
			*s = append(*s, 0)
			copy((*s)[i+1:], (*s)[i:])
			(*s)[i] = x
		}
	} else {
		*s = append(*s, x)
	}
	return true
}

func SortedIntSliceRemove(s *[]int, x int) {
	i := SortedIntSliceSearch(*s, x)
	found := (i != len(*s) && (*s)[i] == x)
	if found {
		*s = append((*s)[:i], (*s)[i+1:]...)
	}
}
