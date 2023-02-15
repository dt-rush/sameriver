package sameriver

import (
	"bytes"
	"os"
	"strings"
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
		debugGOAPPrintf("    nil")
		return
	}
	for spec, interval := range g.goals {
		split := strings.Split(spec, ",")
		varName := split[0]
		debugGOAPPrintf("    %s: [%.0f, %.0f]", varName, interval.a, interval.b)
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
