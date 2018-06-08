package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

// generate a file containing SafeGet methods for each component
func generateComponentsSafeGetMethodsFile(
	components []ComponentSpec) *File {

	f := NewFile("engine")

	for _, component := range components {
		constName := fmt.Sprintf("%s_COMPONENT", strings.ToUpper(component.Name))
		methodName := fmt.Sprintf("SafeGet%s", component.Name)
		f.Func().
			Params(Id("em").Op("*").Id("EntityManager")).
			Id(methodName).
			Params(Id("entity").Id("EntityToken")).
			Params(Id(component.Type), Error()).
			Block(
				If(
					Op("!").Id("em").Dot("lockEntityComponent").
						Call(Id("e"), Id(constName)),
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
				Id("val").Op(":=").Do(func(s *Statement) {
					deepCopyMethodName := fmt.Sprintf("DeepCopy%s",
						component.Name)
					val := Id("em").Dot("Components").
						Dot(component.Name).Index(Id("e").Dot("ID"))
					if component.HasDeepCopyMethod {
						s.Id(deepCopyMethodName).Call(val)
					} else {
						s.Add(val)
					}
				}),
			).Line()
	}

	return f
}
