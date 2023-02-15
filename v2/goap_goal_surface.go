package sameriver

type GOAPGoalRemaining struct {
	goal  *GOAPGoal
	diffs map[string]float64
}

type GOAPGoalRemainingSurface struct {
	main *GOAPGoalRemaining
	pres []*GOAPGoalRemaining
}

func NewGOAPGoalRemainingSurface() *GOAPGoalRemainingSurface {
	return &GOAPGoalRemainingSurface{
		main: nil,
		pres: []*GOAPGoalRemaining{},
	}
}
