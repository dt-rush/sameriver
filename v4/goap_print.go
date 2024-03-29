package sameriver

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

func GOAPPathToString(path *GOAPPath) string {
	var buf bytes.Buffer
	buf.WriteString("    [")
	for i, action := range path.path {
		buf.WriteString(action.DisplayName())
		if i != len(path.path)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]    ")
	return buf.String()
}

func debugGOAPGoalToString(g *GOAPGoal) string {
	if g == nil || len(g.vars) == 0 {
		return color.InBlackOverGreen("    satisfied    ")
	}
	msg := ""
	for varName, interval := range g.vars {
		varInterval := fmt.Sprintf("%s: [%.0f, %.0f]", varName, interval.A, interval.B)
		msg = fmt.Sprintf("%s  %s", msg, color.InRedOverWhite(color.InBold(varInterval)))
	}
	return msg
}

func debugGOAPPrintGoalRemaining(g *GOAPGoalRemaining) {
	if g.nUnfulfilled == 0 {
		msg := "    satisfied    "
		logGOAPDebug(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		logGOAPDebug(color.InBlackOverGreen(msg))
		logGOAPDebug(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		return
	}
	for varName, interval := range g.goalLeft {
		msg := fmt.Sprintf("    %s: [%.0f, %.0f]    ", varName, interval.A, interval.B)

		logGOAPDebug(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))
		logGOAPDebug(color.InBold(color.InRedOverBlack(msg)))
		logGOAPDebug(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))

	}
}

func debugGOAPPrintGoalRemainingSurface(s *GOAPGoalRemainingSurface) {
	if s.NUnfulfilled() == 0 {
		logGOAPDebug(color.InYellowOverGreen("    none remaining    "))
	}
	logGOAPDebug(color.InBold(color.InRedOverGray("%d pres:")), len(s.surface)-1)
	for i, tgs := range s.surface {
		logGOAPDebug("surface[%d] (%d region(s)):", i, len(s.surface[i]))
		if i == len(s.surface)-1 {
			logGOAPDebug(color.InBold(color.InRedOverGray("%d main:")), len(tgs))
		}
		for _, tg := range tgs {
			debugGOAPPrintGoalRemaining(tg)
		}
	}
}
