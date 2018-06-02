package build

import (
	"bytes"
	. "github.com/dave/jennifer/jen"
	"strings"
)

func generateComponentsEnumFile(
	componentNames []string) string {

	// for each component name, create an uppercase const name
	constNames := make(map[string]string)
	for _, componentName := range componentNames {
		constNames[componentName] = strings.ToUpper(componentName) + "_COMPONENT"
	}
	// generate the source file
	var buffer bytes.Buffer

	Type().Id("ComponentType").Int().
		Render(&buffer)
	buffer.WriteString("\n\n")

	Const().Id("N_COMPONENT_TYPES").Op("=").Lit(len(componentNames)).
		Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum
	constDefs := make([]Code, len(componentNames))
	for i, componentName := range componentNames {
		constDefs[i] = Id(constNames[componentName]).Op("=").Iota()
	}
	Const().Defs(constDefs...).
		Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum->string function
	Var().Id("COMPONENT_NAMES").Op("=").
		Map(Id("ComponentType")).String().
		Values(DictFunc(func(d Dict) {
			for _, componentName := range componentNames {
				constName := constNames[componentName]
				d[Id(constName)] = Lit(constName)
			}
		})).
		Render(&buffer)
	return buffer.String()
}
