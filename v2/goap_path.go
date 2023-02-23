package sameriver

import (
	"bytes"
)

type GOAPPath struct {
	path        []*GOAPAction
	cost        int                       // set by GOAPPath.inserted()
	statesAlong []*GOAPWorldState         // set in GOAPEvaluator.computeRemainingsOfPath()
	remainings  *GOAPGoalRemainingSurface // set in GOAPEvaluator.computeRemainingsOfPath
}

func NewGOAPPath(path []*GOAPAction) *GOAPPath {
	if path == nil {
		path = make([]*GOAPAction, 0)
	}
	return &GOAPPath{
		path: path,
	}
}

func (p *GOAPPath) costOfAdd(a *GOAPAction) int {
	// compute cost
	cost := p.cost
	switch a.cost.(type) {
	case int:
		cost += a.cost.(int)
	case func() int:
		cost += a.cost.(func() int)()
	}
	return cost
}

/*
// path [A B]
insertionIx: 0
newSlice: [_ _ _]
copy(newSlice[:0], path[:0]) // [_ _ _]
newslice[0] = a              // [X _ _]
copy(newslice[1:], path[0:]) // [X A B]

// path: []
insertionIx: 0
newSlice: [_]
copy(newSlice[:0], path[:0]) // [_]
newslice[0] = a              // [X]
copy(newslice[1:], path[0:]) // [X]
*/
func (p *GOAPPath) inserted(a *GOAPAction, insertionIx int) *GOAPPath {
	// copy actions into new slice
	newSlice := make([]*GOAPAction, len(p.path)+1)
	copy(newSlice[:insertionIx], p.path[:insertionIx])
	newSlice[insertionIx] = a
	copy(newSlice[insertionIx+1:], p.path[insertionIx:])
	result := &GOAPPath{
		path: newSlice,
		cost: p.costOfAdd(a),
	}
	return result
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
