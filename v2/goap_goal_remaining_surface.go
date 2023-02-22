package sameriver

type GOAPGoalRemainingSurface struct {
	main         *GOAPGoalRemaining
	pres         []*GOAPGoalRemaining
	nUnfulfilled int
	path         *GOAPPath
}

func NewGOAPGoalRemainingSurface() *GOAPGoalRemainingSurface {
	return &GOAPGoalRemainingSurface{
		main: nil,
		pres: []*GOAPGoalRemaining{},
	}
}

func (after *GOAPGoalRemainingSurface) isCloser(before *GOAPGoalRemainingSurface) (closer bool) {
	debugGOAPPrintf("      ** is surface closer?")
	if after.nUnfulfilled < before.nUnfulfilled {
		debugGOAPPrintf("      ** nUnfulfilled was less (after: %d, before: %d)", after.nUnfulfilled, before.nUnfulfilled)
		return true
	}
	if after.main.isCloser(before.main) {
		debugGOAPPrintf("      ** main goal was closer")
		return true
	}
	switch after.path.construction {
	case GOAP_PATH_PREPEND:
		for i := 1; i < len(after.pres); i++ {
			if after.pres[i].isCloser(before.pres[i-1]) {
				debugGOAPPrintf("      ** pre for %s was closer", after.path.path[i].DisplayName())
				return true
			}
		}
	case GOAP_PATH_APPEND:
		for i := 0; i < len(before.pres); i++ {
			if after.pres[i].isCloser(before.pres[i]) {
				debugGOAPPrintf("      ** pre for %s was closer", after.path.path[i].DisplayName())
				return true
			}
		}
	}
	return false
}
