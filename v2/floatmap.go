package sameriver

type FloatMap struct {
	m map[string]float64
}

func NewFloatMap(m map[string]float64) FloatMap {
	return FloatMap{m}
}

func (m *FloatMap) CopyOf() FloatMap {
	m2 := make(map[string]float64)
	for key := range m.m {
		m2[key] = m.m[key]
	}
	return NewFloatMap(m2)
}
