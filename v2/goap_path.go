package sameriver

type GOAPPath struct {
	path []*GOAPAction
}

func NewGOAPPath(path []*GOAPAction) *GOAPPath {
	return &GOAPPath{
		path: path,
	}
}

func (p *GOAPPath) prepend(a *GOAPAction) *GOAPPath {
	newPath := make([]*GOAPAction, len(p.path)+1)
	copy(newPath[1:], p.path)
	newPath[0] = a
	result := &GOAPPath{
		path: newPath,
	}
	return result
}
