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
	// region offsets[i][j] tells where we should insert an action to satisfy
	// surface[i][j]
	regionOffsets [][]int
}

func NewGOAPGoalRemainingSurface(length int) *GOAPGoalRemainingSurface {
	s := &GOAPGoalRemainingSurface{
		surface: make([][]*GOAPGoalRemaining, length),
	}
	for i := range s.surface {
		s.surface[i] = make([]*GOAPGoalRemaining, 0)
	}
	// when this is called for the root goal, the caller wants to be able
	// to set regionOffsets[0] = make([]int, len(main.temporalGoals))
	// but when this is called for the purposes of computeremaining,
	// the caller will just disregard this default value and
	// set the value of s.regionOffsets to the result of calling
	// newRegionOffsetsAfterInsert(i, regionIx)
	s.regionOffsets = [][]int{[]int{}}
	return s
}

func (s *GOAPGoalRemainingSurface) newRegionOffsetsAfterInsert(i int, regionIx int) [][]int {
	newOffsets := make([][]int, len(s.regionOffsets))
	for i, slice := range s.regionOffsets {
		newOffsets[i] = make([]int, len(slice))
		copy(newOffsets[i], slice)
	}
	// TODO: magic here
	return newOffsets
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
