package sameriver

import (
	"math"
)

type GOAPGoalRemaining struct {
	goal  *GOAPGoal
	diffs map[string]float64
}

func (after *GOAPGoalRemaining) isCloser(before *GOAPGoalRemaining) (less bool) {
	debugGOAPPrintf("        *** is remaining less?")
	debugGOAPPrintf("            after.diffs: %v", after.diffs)
	for varName, diff := range after.diffs {
		if math.Abs(diff) < math.Abs(before.diffs[varName]) {
			debugGOAPPrintf("        *** diff for %s was less!", varName)
			return true
		}
	}
	debugGOAPPrintf("        *** not")
	return false
}
