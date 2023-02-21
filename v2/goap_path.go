package sameriver

const (
	GOAP_PATH_PREPEND = iota
	GOAP_PATH_APPEND  = iota
)

type GOAPPath struct {
	path         []*GOAPAction
	cost         int
	construction int
	remainings   *GOAPGoalRemainingSurface
	endState     *GOAPWorldState
}

func NewGOAPPath(path []*GOAPAction, construction int) *GOAPPath {
	return &GOAPPath{
		path:         path,
		construction: construction,
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

func (p *GOAPPath) prepended(a *GOAPAction) *GOAPPath {
	// copy actions into new slice
	newSlice := make([]*GOAPAction, len(p.path)+1)
	copy(newSlice[1:], p.path)
	newSlice[0] = a
	result := &GOAPPath{
		path:         newSlice,
		construction: GOAP_PATH_PREPEND,
		cost:         p.costOfAdd(a),
	}

	return result
}

func (p *GOAPPath) appended(a *GOAPAction) *GOAPPath {
	newPath := make([]*GOAPAction, len(p.path)+1)
	copy(newPath, p.path)
	newPath[len(newPath)-1] = a
	result := &GOAPPath{
		path:         newPath,
		construction: GOAP_PATH_APPEND,
		cost:         p.costOfAdd(a),
	}
	return result
}
