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

func removeEntityTokenFromSlice(slice *[]*EntityToken, x *EntityToken) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeAtIndexInEntityTokenSlice(slice *[]*EntityToken, index int) {
	last_ix := len(*slice) - 1
	(*slice)[index] = (*slice)[last_ix]
	*slice = (*slice)[:last_ix]
}

func indexOfEntityTokenInSlice(slice *[]*EntityToken, x *EntityToken) int {
	for i, v := range *slice {
		if v == x {
			return i
		}
	}
	return -1
}

func appendStringToSlice(slice *[]string, x string) {
	*slice = append(*slice, x)
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

func removeEventChannelFromSlice(slice *[]*EventChannel, x *EventChannel) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v.C == x.C {
			(*slice)[i] = (*slice)[last_ix]
			// set the last element to the zero value (for same reasons
			// as above)
			(*slice)[last_ix] = &EventChannel{}
			*slice = (*slice)[:last_ix]
			break
		}
	}
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
