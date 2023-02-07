package sameriver

import (
	"container/heap"
)

type GOAPPQueueItem struct {
	path []*GOAPAction
	goal *GOAPGoal
	cost int
	// The index is needed by update and is maintained
	// by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func NewGOAPPQueueItem(path []*GOAPAction, goal *GOAPGoal) *GOAPPQueueItem {
	// add path cost
	cost := 0
	for _, action := range path {
		switch action.cost.(type) {
		case int:
			cost += action.cost.(int)
		case func() int:
			cost += action.cost.(func() int)()
		}
	}
	// add heuristic (number of unfulfilled state vars remaining)
	cost += len(goal.goals)
	return &GOAPPQueueItem{
		path,
		goal,
		cost,
		-1, // going to be set on Push() anyway
	}
}

type GOAPPriorityQueue []*GOAPPQueueItem

func (pq GOAPPriorityQueue) Len() int { return len(pq) }

func (pq GOAPPriorityQueue) Less(i, j int) bool {
	// We goal Pop to give us the lowest cost
	return pq[i].cost < pq[j].cost
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

// update modifies the cost and value of an
// Item in the queue.
func (pq *GOAPPriorityQueue) update(
	item *GOAPPQueueItem,
	path []*GOAPAction,
	cost int) {

	item.path = path
	item.cost = cost
	heap.Fix(pq, item.index)
}
