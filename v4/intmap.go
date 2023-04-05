package sameriver

type IntMap struct {
	m              map[string]int
	validIntervals map[string][2]int
}

func NewIntMap(m map[string]int) IntMap {
	return IntMap{m, make(map[string][2]int)}
}

func (m *IntMap) CopyOf() IntMap {
	m2 := make(map[string]int)
	for key := range m.m {
		m2[key] = m.m[key]
	}
	return IntMap{m2, m.validIntervals}
}

func (m *IntMap) SetValidInterval(k string, a, b int) {
	m.validIntervals[k] = [2]int{a, b}
}

func (m *IntMap) ValCanBeSetTo(k string, v int) bool {
	if validInterval, exists := m.validIntervals[k]; exists {
		return v >= validInterval[0] && v <= validInterval[1]
	} else {
		return true
	}
}

func (m *IntMap) Set(k string, v int) {
	if validInterval, exists := m.validIntervals[k]; exists {
		if v < validInterval[0] {
			v = validInterval[0]
		} else if v > validInterval[1] {
			v = validInterval[1]
		}
	}
	m.m[k] = v
}

func (m *IntMap) Get(k string) int {
	return m.m[k]
}
