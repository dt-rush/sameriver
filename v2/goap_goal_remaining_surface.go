package sameriver

import (
	"github.com/dt-rush/sameriver/v2/utils"
)

type GOAPGoalRemaining struct {
	goal         *GOAPGoal
	goalLeft     map[string]*utils.NumericInterval
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
