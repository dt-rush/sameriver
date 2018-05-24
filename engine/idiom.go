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

package engine

// thanks to https://stackoverflow.com/a/37359662 for this nice
// little splice idiom when we don't care about slice order (saves
// a copy operation if we wanted to shift the slice to fill the gap)

func removeUint16FromSlice(slice *[]uint16, x uint16) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeIndexFromUint16Slice(slice *[]uint16, index int) {
	last_ix := len(*slice) - 1
	(*slice)[index] = (*slice)[last_ix]
	*slice = (*slice)[:last_ix]
}

func indexOfUint16InSlice(slice *[]uint16, x uint16) int {
	for i, v := range *slice {
		if v == x {
			return i
		}
	}
	return -1
}

func removeStringFromSlice(slice *[]string, x string) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeEntityTokenFromSlice(slice *[]EntityToken, x EntityToken) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v.ID == x.ID {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeEntityQueryWatcherFromSliceByID(
	slice *[]EntityQueryWatcher, ID uint16) {

	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v.ID == ID {
			(*slice)[i] = (*slice)[last_ix]
			// set the last element (which we will then cut off the end)
			// to the zero value, so that we don't leave any pointer members
			// still sitting there in the shadow of the slice backing array
			(*slice)[last_ix] = EntityQueryWatcher{}
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeEventChannelFromSlice(
	slice *[]EventChannel, x EventChannel) {

	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v.C == x.C {
			(*slice)[i] = (*slice)[last_ix]
			// set the last element to the zero value (for same reasons
			// as above)
			(*slice)[last_ix] = EventChannel{}
			*slice = (*slice)[:last_ix]
			break
		}
	}
}
