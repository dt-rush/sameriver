package generate

import (
	. "github.com/dave/jennifer/jen"
)

func generateComponentsTableFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	// build the ComponentsTable struct declaration
	f.Type().Id("ComponentsTable").StructFunc(func(g *Group) {
		g.Id("em").Op("*").Id("EntityManager")
		for _, component := range components {
			g.Id(component.Name).
				Index(Id("MAX_ENTITIES")).Id(component.Type)
		}
	}).Line()

	// write the Init method
	f.Func().
		Id("NewComponentsTable").
		Params(Id("em").Op("*").Id("EntityManager")).
		Op("*").Id("ComponentsTable").
		Block(
			Return(Op("&").Id("ComponentsTable").
				Values(DictFunc(func(d Dict) {
					d[Id("em")] = Id("em")
				}))),
		).Line()

	return f
}
