package sameriver

type GOAPGoalRemaining struct {
	goal     *GOAPGoal
	goalLeft map[string]*NumericInterval

	diffs        map[string]float64
	nUnfulfilled int
}

type GOAPGoalRemainingSurface struct {
	// for each action in the path, there is a []*GOAPGoalRemaining,
	// representing the temporal goal remainings
	// so for example if we have path [A, B, C], end goal u,
	// where C preconditions are [s t],
	// A preconditions are [q] and
	// B preconditions are [r],
	// then our surface will be
	// [ [q] [r] [s t] [u] ]
	surface [][]*GOAPGoalRemaining
}

func newGOAPGoalRemainingSurface(length int) *GOAPGoalRemainingSurface {
	s := &GOAPGoalRemainingSurface{
		surface: make([][]*GOAPGoalRemaining, length),
	}
	for i := range s.surface {
		s.surface[i] = make([]*GOAPGoalRemaining, 0)
	}

	return s
}

func (s *GOAPGoalRemainingSurface) NUnfulfilled() int {
	n := 0
	// for each []*GOAPGoalRemaining of surface
	for i := range s.surface {
		n += s.nUnfulfilledAtIx(i)
	}
	return n
}

func (s *GOAPGoalRemainingSurface) nUnfulfilledAtIx(i int) int {
	n := 0
	for _, tg := range s.surface[i] {
		n += tg.nUnfulfilled
	}
	return n
}
