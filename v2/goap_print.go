package sameriver

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

var DEBUG_GOAP_VAL, DEBUG_GOAP_OK = os.LookupEnv("DEBUG_GOAP")
var DEBUG_GOAP = DEBUG_GOAP_OK && DEBUG_GOAP_VAL == "true"

func debugGOAPPrintf(s string, args ...any) {
	if DEBUG_GOAP {
		Logger.Printf(s, args...)
	}
}

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
		msg = color.InRedOverBlack(fmt.Sprintf("%s  %s", msg, varInterval))
	}
	return msg
}

func debugGOAPPrintGoalRemaining(g *GOAPGoalRemaining) {
	if g.nUnfulfilled == 0 {
		msg := "    satisfied    "
		debugGOAPPrintf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		debugGOAPPrintf(color.InBlackOverGreen(msg))
		debugGOAPPrintf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		return
	}
	for varName, interval := range g.goalLeft {
		msg := fmt.Sprintf("    %s: [%.0f, %.0f]    ", varName, interval.A, interval.B)

		debugGOAPPrintf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))
		debugGOAPPrintf(color.InBold(color.InRedOverBlack(msg)))
		debugGOAPPrintf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))

	}
}

func debugGOAPPrintGoalRemainingSurface(s *GOAPGoalRemainingSurface) {
	if s.NUnfulfilled() == 0 {
		debugGOAPPrintf(color.InYellowOverGreen("    none remaining    "))
	}
	debugGOAPPrintf(color.InBold(color.InRedOverGray("pres:")))
	for i, goal := range s.surface {
		if i == len(s.surface)-1 {
			debugGOAPPrintf(color.InBold(color.InRedOverGray("main:")))
			debugGOAPPrintGoalRemaining(s.surface[len(s.surface)-1])
		} else {
			debugGOAPPrintGoalRemaining(goal)
		}
	}
}

func debugGOAPPrintWorldState(ws *GOAPWorldState) {
	if ws == nil || len(ws.vals) == 0 {
		debugGOAPPrintf("    nil")
		return
	}
	for name, val := range ws.vals {
		debugGOAPPrintf("    %s: %d", name, val)
	}
}
