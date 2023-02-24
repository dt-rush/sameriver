package sameriver

type GOAPGoalRemaining struct {
	goal         *GOAPGoal
	goalLeft     map[string]*NumericInterval
	diffs        map[string]float64
	nUnfulfilled int
}

type GOAPGoalRemainingSurface struct {
	surface []*GOAPGoalRemaining
}

func NewGOAPGoalRemainingSurface() *GOAPGoalRemainingSurface {
	return &GOAPGoalRemainingSurface{
		surface: []*GOAPGoalRemaining{},
	}
}

func (s *GOAPGoalRemainingSurface) NUnfulfilled() int {
	n := 0
	for _, g := range s.surface {
		n += g.nUnfulfilled
	}
	return n
}
