package main

import (
	"bytes"
	"errors"
	"github.com/disiqueira/gotree"
)

type NodeHeap struct {
	Arr []*Node
}

func NewNodeHeap() *NodeHeap {
	// leaving first position unused so bit shifting works well
	h := NodeHeap{Arr: []*Node{nil}}
	return &h
}

func (h *NodeHeap) bubbleUp(ix int) int {
	for ix > 1 && h.Arr[ix].F < h.Arr[ix>>1].F {
		h.Arr[ix].HeapIX, h.Arr[ix>>1].HeapIX = h.Arr[ix>>1].HeapIX, h.Arr[ix].HeapIX
		h.Arr[ix], h.Arr[ix>>1] = h.Arr[ix>>1], h.Arr[ix]
		ix = ix >> 1
	}
	return ix
}

func (h *NodeHeap) Add(x *Node) (ix int) {
	x.F = x.G + x.H
	h.Arr = append(h.Arr, x)
	ix = len(h.Arr) - 1
	h.Arr[ix].HeapIX = ix
	ix = h.bubbleUp(ix)
	h.Arr[ix].HeapIX = ix
	return ix
}

func (h *NodeHeap) bubbleDown(ix int) int {
	for {
		greater := ix
		lix := (ix << 1)
		rix := (ix << 1) + 1
		// if we've reached the bottom, return
		if !(lix < len(h.Arr) || rix < len(h.Arr)) {
			return ix
		}
		// check if left node is greater
		if lix < len(h.Arr) && h.Arr[lix].F < h.Arr[greater].F {
			greater = lix
		}
		// check if right node is greater
		if rix < len(h.Arr) && h.Arr[rix].F < h.Arr[greater].F {
			greater = rix
		}
		// if one of children was greater, swap and continue bubble down
		// from that node
		if greater != ix {
			h.Arr[ix].HeapIX, h.Arr[greater].HeapIX = h.Arr[greater].HeapIX, h.Arr[ix].HeapIX
			h.Arr[ix], h.Arr[greater] = h.Arr[greater], h.Arr[ix]
			h.bubbleDown(greater)
			ix = greater
		} else {
			// else, no child was greater. we can stop here
			return ix
		}
	}
}

func (h *NodeHeap) Pop() (*Node, error) {
	if h.Len() == 0 {
		return nil, errors.New("heap empty")
	}
	// get root elem and replace with last element (shrink slice)
	x := h.Arr[1]
	last_ix := len(h.Arr) - 1
	h.Arr[1] = h.Arr[last_ix]
	h.Arr[1].HeapIX = 1
	h.Arr = h.Arr[:last_ix]
	// bubble element down to its place
	h.bubbleDown(1)
	return x, nil
}

func (h *NodeHeap) Modify(ix int, G int) {
	// get the old F value
	oldVal := h.Arr[ix].F
	// set the new G value and calculate the new F value
	h.Arr[ix].G = G
	h.Arr[ix].F = h.Arr[ix].G + h.Arr[ix].H
	// we will need to modify the node's HeapIX value once it bubbles
	var newIX = ix
	// bubble up if needed
	if h.Arr[ix].F < oldVal {
		newIX = h.bubbleUp(ix)
	}
	// bubble down if needed
	if h.Arr[ix].F > oldVal {
		newIX = h.bubbleDown(ix)
	}
	h.Arr[newIX].HeapIX = newIX
}

func (h *NodeHeap) Len() int {
	return len(h.Arr) - 1
}

func (h *NodeHeap) Clear() {
	h.Arr = h.Arr[:1]
}

func (h *NodeHeap) String() string {
	// for building string
	var buffer bytes.Buffer
	// print array
	// if elements, print tree
	if len(h.Arr) > 1 {
		buffer.WriteString("\n")
		// build tree using gotree package by descending recursively
		var addChildren func(node gotree.Tree, ix int)
		addChildren = func(node gotree.Tree, ix int) {
			rix := (ix << 1) + 1
			if rix < len(h.Arr) {
				r := node.Add(h.Arr[rix].String())
				addChildren(r, rix)
			}
			lix := (ix << 1)
			if lix < len(h.Arr) {
				l := node.Add(h.Arr[lix].String())
				addChildren(l, lix)
			}
		}
		tree := gotree.New(h.Arr[1].String())
		addChildren(tree, 1)
		buffer.WriteString(tree.Print())
	}
	return buffer.String()
}
