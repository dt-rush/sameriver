package sameriver

import (
	"bytes"
)

type GOAPPath struct {
	path []*GOAPAction
	cost int // set by GOAPPath.inserted()
	// states after each action, from start state at [0]
	// til the end state after the last action
	statesAlong []*GOAPWorldState // set in GOAPEvaluator.computeRemainingsOfPath()
	// goal remaining for each action's pres in the path
	// plus the goal remaining for the main goal at
	// the last index
	remainings *GOAPGoalRemainingSurface // set in GOAPEvaluator.computeRemainingsOfPath
	// parallel array to path:
	// the region offsets to insert the temporal groups of the pres for each
	// action at i
	regionOffsets [][]int
}

func NewGOAPPath(path []*GOAPAction) *GOAPPath {
	if path == nil {
		path = make([]*GOAPAction, 0)
	}
	return &GOAPPath{
		path:          path,
		regionOffsets: [][]int{[]int{}},
	}
}

func (p *GOAPPath) costOfAdd(a *GOAPAction) int {
	// compute cost
	cost := p.cost
	switch a.cost.(type) {
	case int:
		cost += a.Count * a.cost.(int)
	case func() int:
		cost += a.Count * a.cost.(func() int)()
	}
	return cost
}

// regionIx is the region index this action was inserted into, satisfying
func (p *GOAPPath) inserted(a *GOAPAction, insertionIx int, regionIx int) *GOAPPath {
	// copy actions into new slice, and put a at insertionIx
	newSlice := make([]*GOAPAction, len(p.path)+1)
	for i := 0; i < insertionIx; i++ {
		newSlice[i] = p.path[i].CopyOf()
	}
	newSlice[insertionIx] = a
	for i := insertionIx + 1; i < len(newSlice); i++ {
		newSlice[i] = p.path[i-1].CopyOf()
	}

	path := &GOAPPath{
		path: newSlice,
		cost: p.costOfAdd(a),
	}
	a.insertionIx = insertionIx
	a.regionIx = regionIx
	// copy regionOffsets data
	path.regionOffsets = make([][]int, len(p.regionOffsets))
	for i := 0; i < len(p.regionOffsets); i++ {
		path.regionOffsets[i] = make([]int, len(p.regionOffsets[i]))
		copy(path.regionOffsets[i], p.regionOffsets[i])
	}

	// go up tree, updating regionOffsets
	node := a
	for node != nil {
		// if goal isnt' temporal (length 1), no update at this level
		// (note remainings here is before including the pres of a,
		// and we haven't udpated the insertionIx of the old path actions,
		// so we can use node.insertionIx of a node (from a.parent on up)
		// to get its goal surface
		// surfaceIx is the surface index of the temporal goal this action satisfies
		var surfaceIx int
		if node.parent == nil {
			surfaceIx = len(p.remainings.surface) - 1
		} else {
			logGOAPDebug("    node.parent = %s, parent.insertionIx = %d", node.parent.Name, node.parent.insertionIx)
			surfaceIx = node.parent.insertionIx
		}
		logGOAPDebug("    insert.parent surfaceIx: %d", surfaceIx)
		if len(p.remainings.surface[surfaceIx]) == 1 {
			node = node.parent
			continue
		}
		// if we inserted to regionIx 0, there's nothing to the left of it needed an updated offset
		if node.regionIx == 0 {
			node = node.parent
			continue
		}
		// thus, regionIx > 1 had the insertion, and we need to shift the regions to its left
		// by -1
		for ri := node.regionIx - 1; ri >= 0; ri-- {
			logGOAPDebug("            updating regionOffsets[%d][%d]", surfaceIx, ri)
			path.regionOffsets[surfaceIx][ri] -= 1
		}
		node = node.parent
	}
	// insert regionOffsets for a.pres
	path.regionOffsets = append(
		path.regionOffsets[:insertionIx+1],
		path.regionOffsets[insertionIx:]...)
	path.regionOffsets[insertionIx] = make([]int, len(a.pres.temporalGoals))
	if DEBUG_GOAP {
		logGOAPDebug("  regionOffsets after insert&update: %v", path.regionOffsets)
	}
	// update action indexes after insertion
	for j := insertionIx + 1; j < len(path.path); j++ {
		path.path[j].insertionIx++
	}
	insertionIxs := make([]int, len(path.path))
	for i, ac := range path.path {
		insertionIxs[i] = ac.insertionIx
	}
	logGOAPDebug("  insertionIxs after insert&update: %v", insertionIxs)
	// add regionOffsets
	// update this aciton's index
	return path
}

func (p *GOAPPath) String() string {
	var buf bytes.Buffer
	buf.WriteString("    [")
	for i, action := range p.path {
		buf.WriteString(action.DisplayName())
		if i != len(p.path)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]    ")
	return buf.String()
}
