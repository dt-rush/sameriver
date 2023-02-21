package sameriver

type GOAPPQueueItem struct {
	path *GOAPPath
	// The index is needed by update and is maintained
	// by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type GOAPPriorityQueue []*GOAPPQueueItem

func (pq GOAPPriorityQueue) Len() int { return len(pq) }

func (pq GOAPPriorityQueue) Less(i, j int) bool {
	// We goal Pop to give us the lowest cost
	return pq[i].path.cost < pq[j].path.cost
}

func (pq GOAPPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *GOAPPriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*GOAPPQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *GOAPPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
