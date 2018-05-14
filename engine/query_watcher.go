package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type QueryWatcher struct {
	Query   bitarray.BitArray
	Channel chan (int)
	ID      int
}
