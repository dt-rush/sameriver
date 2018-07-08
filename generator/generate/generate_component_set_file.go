package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

func generateComponentSetFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	// create struct definition
	f.Type().Id("ComponentSet").StructFunc(func(g *Group) {
		for _, component := range components {
			g.Id(component.Name).Op("*").Id(component.Type)
		}
	}).Line()

	// create ToBitArray() method
	bitArrayPkg := "github.com/golang-collections/go-datastructures/bitarray"
	f.ImportName(bitArrayPkg, "bitarray")
	f.Func().
		Params(Id("cs").Op("*").Id("ComponentSet")).
		Id("ToBitArray").Params().
		Qual(bitArrayPkg, "BitArray").
		BlockFunc(func(g *Group) {
			g.Id("b").Op(":=").Qual(bitArrayPkg, "NewBitArray").
				Call(Uint64().Parens(Id("N_COMPONENT_TYPES")))
			for _, component := range components {
				constName := strings.ToUpper(
					fmt.Sprintf("%s_COMPONENT", component.Name))
				g.If(Id("cs").Dot(component.Name).Op("!=").Nil()).
					Block(
						Id("b").Dot("SetBit").Call(
							Uint64().Parens(Id(constName))),
					)
			}
			g.Return(Id("b"))
		}).Line()

	// create EntityManager.ApplyComponent method
	f.Func().
		Params(Id("em").Op("*").Id("EntityManager")).
		Id("ApplyComponentSet").Params(Id("cs").Id("ComponentSet")).
		Func().Parens(Id("*EntityToken")).
		Block(
			Return(Func().Parens(Id("entity").Op("*").Id("EntityToken")).
				BlockFunc(func(g *Group) {
					for _, component := range components {
						g.If(Id("cs").Dot(component.Name).Op("!=").Nil()).Block(
							Id("em").Dot("Components").Dot(component.Name).
								Index(Id("entity").Dot("ID")).Op("=").
								Op("*").Id("cs").Dot(component.Name),
						)
					}
				}),
			),
		)

	return f
}
