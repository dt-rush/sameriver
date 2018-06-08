package generate

import (
	. "github.com/dave/jennifer/jen"
	"strings"
)

func generateComponentsEnumFile(
	components []ComponentSpec) *File {

	// for each component name, create an uppercase const name
	constNames := make(map[string]string)
	for _, component := range components {
		constNames[component.Name] =
			strings.ToUpper(component.Name) + "_COMPONENT"
	}

	f := NewFile("engine")

	f.Type().Id("ComponentType").Int()

	f.Const().Id("N_COMPONENT_TYPES").Op("=").Lit(len(components))

	// write the enum
	f.Const().DefsFunc(func(g *Group) {
		for _, component := range components {
			g.Id(constNames[component.Name]).Op("=").Iota()
		}
	})

	// write the enum->string function
	f.Var().Id("COMPONENT_NAMES").Op("=").
		Map(Id("ComponentType")).String().
		Values(DictFunc(func(d Dict) {
			for _, component := range components {
				constName := constNames[component.Name]
				d[Id(constName)] = Lit(constName)
			}
		}))

	return f
}
