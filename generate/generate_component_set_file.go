package generate

import (
	. "github.com/dave/jennifer/jen"
)

func generateComponentSetFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	f.Type().Id("ComponentSet").StructFunc(func(g *Group) {
		for _, component := range components {
			g.Id(component.Name).Op("*").Id(component.Type)
		}
	})

	return f
}
