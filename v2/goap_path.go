package sameriver

const (
	GOAP_PATH_PREPEND = iota
	GOAP_PATH_APPEND  = iota
)

type GOAPPath struct {
	path         []*GOAPAction
	construction int
}

func NewGOAPPath(path []*GOAPAction, construction int) *GOAPPath {
	return &GOAPPath{
		path:         path,
		construction: construction,
	}
}

func (p *GOAPPath) prepend(a *GOAPAction) *GOAPPath {
	newPath := make([]*GOAPAction, len(p.path)+1)
	copy(newPath[1:], p.path)
	newPath[0] = a
	result := &GOAPPath{
		path:         newPath,
		construction: GOAP_PATH_PREPEND,
	}
	return result
}

func (p *GOAPPath) append(a *GOAPAction) *GOAPPath {
	newPath := make([]*GOAPAction, len(p.path)+1)
	copy(newPath, p.path)
	newPath[len(newPath)-1] = a
	result := &GOAPPath{
		path:         newPath,
		construction: GOAP_PATH_APPEND,
	}
	return result
}
