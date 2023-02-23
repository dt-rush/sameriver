package sameriver

type IntMap struct {
	m map[string]int
}

func NewIntMap(m map[string]int) IntMap {
	return IntMap{m}
}

func (m *IntMap) CopyOf() IntMap {
	m2 := make(map[string]int)
	for key := range m.m {
		m2[key] = m.m[key]
	}
	return NewIntMap(m2)
}
