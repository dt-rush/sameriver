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
		f.Commentf("read an entity's %s by read-locking the value",
			component.Name)
		f.Comment("this is safe at any time, in any logic goroutine")
		f.Func().
			Params(Id("ct").Op("*").Id("ComponentsTable")).
			Id(methodName).
			Params(Id("e").Id("EntityToken")).
			Params(Id(component.Type), Error()).
			Block(
				Comment("lock the entity's component value as a reader, and"),
				Comment("check the gen is valid once we have the lock."),
				Comment("(release if not valid, returning error)"),
				Id("ct").Dot("em").Dot("rLockEntityComponent").
					Call(Id("e"), Id(constName)),
				Defer().Id("ct").Dot("em").Dot("rUnlockEntityComponent").
					Call(Id("e"), Id(constName)),
				If(
					Op("!").Id("ct").Dot("em").Dot("entityTable").
						Dot("genValidate").Call(Id("e")),
				).Block(
					Comment("if genValidate failed, return an error"),
					Comment("we use the first element of the component data"),
					Comment("simply because it's of the right type."),
					Return(
						List(
							Id("ct").Dot(component.Name).Index(Lit(0)),
							Qual("errors", "New").Call(
								Qual("fmt", "Sprintf").Call(
									Lit("%+v no longer exists"),
									Id("e"),
								),
							),
						),
					),
				),
				Do(func(s *Statement) {
					if component.HasDeepCopyMethod {
						s.Commentf(
							"get the %s with deepcopy, since DeepCopy%s() "+
								"was specified.\n"+
								"(a DeepCopy method should be specified "+
								"whenever the component value\n"+
								"returned may change as a receiver is reading "+
								"it, after it's been\n"+
								"retrieved from the component data table)",
							component.Name, component.Name)
					}
				}),
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
