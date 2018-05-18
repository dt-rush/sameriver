//
// Yes, yes, yes, certain "idioms" of Go are pretty, but some are a little
// awkward to look at, such as removing an element from a slice (without
// preserving order).
//
// This file is made even uglier by the fact that there are no generics...
// just saying.
//

package engine

// thanks to https://stackoverflow.com/a/37359662 for this nice
// little splice idiom when we don't care about slice order (saves
// a copy operation if we wanted to shift the slice to fill the gap)

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

func removeStringFromSlice(x string, slice *[]string) {
	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeEntityQueryWatcherFromSlice(
	x EntityQueryWatcher, slice *[]EntityQueryWatcher) {

	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			(*slice)[last_ix] = nil
			*slice = (*slice)[:last_ix]
			break
		}
	}
}

func removeGameEventChannelFromSlice(
	x GameEventChannel, slice *[]GameEventChannel) {

	last_ix := len(*slice) - 1
	for i, v := range *slice {
		if v == x {
			(*slice)[i] = (*slice)[last_ix]
			(*slice)[last_ix] = nil
			*slice = (*slice)[:last_ix]
			break
		}
	}
}
