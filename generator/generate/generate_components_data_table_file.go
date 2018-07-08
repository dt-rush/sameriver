package generate

import (
	. "github.com/dave/jennifer/jen"
)

func generateComponentsDataTableFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	// build the ComponentsDataTable struct declaration
	f.Type().Id("ComponentsDataTable").StructFunc(func(g *Group) {
		g.Id("em").Op("*").Id("EntityManager")
		for _, component := range components {
			g.Id(component.Name).
				Index(Id("MAX_ENTITIES")).Id(component.Type)
		}
	}).Line()

	// write the Init method
	f.Func().
		Id("NewComponentsDataTable").
		Params(Id("em").Op("*").Id("EntityManager")).
		Op("*").Id("ComponentsDataTable").
		Block(
			Return(Op("&").Id("ComponentsDataTable").
				Values(DictFunc(func(d Dict) {
					d[Id("em")] = Id("em")
				}))),
		).Line()

	return f
}
