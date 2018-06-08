package generate

import (
	. "github.com/dave/jennifer/jen"
)

func generateComponentsTableFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	// build the ComponentsTable struct declaration
	f.Type().Id("ComponentsTable").StructFunc(func(g *Group) {
		g.Id("accessLocks").
			Index(Id("N_COMPONENT_TYPES")).Qual("sync", "RWMutex")
		g.Id("valueLocks").
			Index(Id("N_COMPONENT_TYPES")).
			Index(Id("MAX_ENTITIES")).Qual("sync", "RWMutex")
		for _, component := range components {
			g.Id(component.Name).
				Index(Id("MAX_ENTITIES")).Id(component.Type)
		}
	}).Line()

	// write the Init method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("Init").
		Params(Id("em").Op("*").Id("EntityManager")).
		Block(
			For(
				Id("i").Op(":=").Lit(0),
				Id("i").Op("<").Id("N_COMPONENT_TYPES"),
				Id("i").Op("++"),
			).Block(
				Id("ct").Dot("accessLocks").Index(Id("i")).Op("=").
					Id("NewComponentAccessLock").Call()),
		).Line()

	// write the lock method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("lock").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("Lock").Call(),
		).Line()

	// write the unlock method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("unlock").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("Unlock").Call(),
		).Line()

	// write the accessStart method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("accessStart").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("RLock").Call(),
		).Line()

	// write the accessEnd method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("accessEnd").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("RUnlock").Call(),
		).Line()

	return f
}
