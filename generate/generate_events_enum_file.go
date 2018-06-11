package generate

import (
	. "github.com/dave/jennifer/jen"
	"strings"
)

func generateEventsEnumFile(eventNames []string) *File {
	// for each event name, create an uppercase const name
	constNames := make(map[string]string)
	for _, eventName := range eventNames {
		eventNameStem := strings.Replace(eventName, "Data", "", 1)
		constNames[eventName] = strings.ToUpper(eventNameStem) + "_EVENT"
	}
	// generate the source file
	f := NewFile("engine")

	// typedef EventType int
	f.Type().Id("EventType").Int()

	// const N_EVENT_TYPES = ___
	f.Const().Id("N_EVENT_TYPES").Op("=").Lit(len(eventNames))

	// write the enum
	f.Const().DefsFunc(func(g *Group) {
		for _, eventName := range eventNames {
			g.Id(constNames[eventName]).Op("=").Iota()
		}
	})

	// write the enum->string function
	f.Var().Id("EVENT_NAMES").Op("=").
		Map(Id("EventType")).String().
		Values(DictFunc(func(d Dict) {
			for _, eventName := range eventNames {
				constName := constNames[eventName]
				d[Id(constName)] = Lit(constName)
			}
		}))

	return f

}
