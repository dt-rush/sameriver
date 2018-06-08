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
			Index(Id("N_COMPONENT_TYPES")).Op("*").Id("ComponentAccessLock")
		g.Id("valueLocks").
			Index(Id("N_COMPONENT_TYPES")).
			Index(Id("MAX_ENTITIES")).Qual("sync", "RWMutex")
		for _, component := range components {
			g.Id(component.Name).
				Index(Id("MAX_ENTITIES")).Id(component.Type)
		}
	})

	// write the lock method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("lock").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("Lock").Call(),
		)

	// write the unlock method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("unlock").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("Unlock").Call(),
		)

	// write the access method
	f.Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("access").
		Params(Id("component").Id("ComponentType")).
		Block(
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("RLock").Call(),
			Id("ct").Dot("accessLocks").Index(Id("component")).
				Dot("RUnlock").Call(),
		)

	return f
}
