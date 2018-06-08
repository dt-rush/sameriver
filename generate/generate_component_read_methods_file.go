package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

// generate a file containing Read methods for each component
func generateComponentReadMethodsFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	for _, component := range components {
		constName := fmt.Sprintf("%s_COMPONENT", strings.ToUpper(component.Name))
		methodName := fmt.Sprintf("Read%s", component.Name)
		f.Func().
			Params(Id("ct").Op("*").Id("ComponentsTable")).
			Id(methodName).
			Params(Id("entity").Id("EntityToken")).
			Params(Id(component.Type), Error()).
			Block(
				Id("ct").Dot("em").Dot("rLockEntityComponent").
					Call(Id("e"), Id(constName)),
				Defer().Id("ct").Dot("em").Dot("rUnlockEntityComponent").
					Call(Id("e"), Id(constName)),
				If(
					Op("!").Id("ct").Dot("em").Dot("entityTable").
						Dot("genValidate").Call(Id("e")),
				).Block(
					Return(
						List(
							Id(component.Type).Values(),
							Qual("errors", "New").Dot("New").Call(
								Qual("fmt", "Sprintf").Call(
									Lit("%+v no longer exists"),
									Id("e"),
								),
							),
						),
					),
				),
				Return(
					List(
						Do(func(s *Statement) {
							deepCopyMethodName := fmt.Sprintf("DeepCopy%s",
								component.Name)
							val := Id("ct").Dot(component.Name).
								Index(Id("e").Dot("ID"))
							if component.HasDeepCopyMethod {
								s.Id(deepCopyMethodName).Call(val)
							} else {
								s.Add(val)
							}
						}),
						Nil(),
					),
				),
			).Line()
	}

	return f
}
