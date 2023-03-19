//
// Yes, yes, yes, certain "idioms" of Go are pretty, but some are a little
// awkward to look at, such as removing an element from a slice (without
// preserving order), as documented in the "Slice Tricks" page of the golang/go
// github wiki:
//
// https://github.com/golang/go/wiki/SliceTricks
//
// This file is made even uglier by the fact that there are no generics...
// A simple compile-time check for static type templating should be fine, no
// need to allow template metaprogramming crust to accumulate in this fine
// and beautiful language
//

package sameriver

func IntMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func IntMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func IntAbs(x int) int {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

// thanks to https://stackoverflow.com/a/37359662 for this nice
// little splice idiom when we don't care about slice order (saves
// a copy operation if we wanted to shift the slice to fill the gap)

func removeEntityFromSlice(slice *[]*Entity, x *Entity) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

// TODO: reverse arguments and remove pointer (we're not modifying!)
func indexOfEntityInSlice(slice *[]*Entity, x *Entity) int {
	for i, v := range *slice {
		if v == x {
			return i
		}
	}
	return -1
}

func removeEventChannelFromSlice(slice []*EventChannel, x *EventChannel) []*EventChannel {
	last_ix := len(slice) - 1
	for i, v := range slice {
		if v.C == x.C {
			slice[i] = slice[last_ix]
			// set the last element to the zero value (for same reasons
			// as above)
			slice[last_ix] = &EventChannel{}
			slice = slice[:last_ix]
			break
		}
	}
	return slice
}

func removeUpdatedEntityListFromSlice(
	slice *[]*UpdatedEntityList, x *UpdatedEntityList) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			// set the last element to the zero value (for same reasons
			// as above)
			(*slice)[last_ix] = nil
			*slice = (*slice)[:last_ix]
			break
		}
	}
}
