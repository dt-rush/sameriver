package sameriver

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

func debugGOAPPrintf(s string, args ...any) {
	if val, ok := os.LookupEnv("DEBUG_GOAP"); ok && val == "true" {
		Logger.Printf(s, args...)
	}
}

func GOAPPathToString(path *GOAPPath) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, action := range path.path {
		buf.WriteString(action.name)
		if i != len(path.path)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func debugGOAPPrintGoal(g *GOAPGoal) {
	if g == nil || len(g.goals) == 0 {
		msg := "    satisfied    "
		debugGOAPPrintf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		debugGOAPPrintf(color.InBlackOverGreen(msg))
		debugGOAPPrintf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		return
	}
	for spec, interval := range g.goals {
		split := strings.Split(spec, ",")
		varName := split[0]
		msg := fmt.Sprintf("    %s: [%.0f, %.0f]    ", varName, interval.A, interval.B)

		debugGOAPPrintf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))
		debugGOAPPrintf(color.InBold(color.InRedOverBlack(msg)))
		debugGOAPPrintf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))

	}
}

func debugGOAPPrintGoalRemainingSurface(g *GOAPGoalRemainingSurface) {
	debugGOAPPrintf(color.InBold(color.InRedOverGray("main:")))
	debugGOAPPrintGoal(g.main.goal)
	debugGOAPPrintf(color.InBold(color.InRedOverGray("pres:")))
	for _, pre := range g.pres {
		debugGOAPPrintGoal(pre.goal)
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
