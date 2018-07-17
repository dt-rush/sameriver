package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
)

func generateEntityComponentsGetMethodsFile(components []ComponentSpec) *File {

	f := NewFile("engine")

	// generate Get methods
	for _, component := range components {
		f.Func().
			Params(Id("e").Op("*").Id("Entity")).
			Id(fmt.Sprintf("Get%s", component.Name)).Params().
			Op("*").Id(component.Type).
			Block(
				Return(
					Op("&").Id("e").Dot("World").Dot("em").Dot("components").
						Dot(component.Name).Index(Id("e").Dot("ID")),
				),
			)
	}

	return f
}
