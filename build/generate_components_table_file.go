package build

import (
	"bytes"
	. "github.com/dave/jennifer/jen"
)

func generateComponentsTableFile(
	componentNames []string) string {

	// generate the source file
	var buffer bytes.Buffer
	// build the ComponentsTable struct declaration
	fields := make([]Code, len(componentNames))
	for i, componentName := range componentNames {
		fields[i] = Id(componentName).
			Op("*").Id(componentStructName(componentName))
	}
	Type().Id("ComponentsTable").Struct(fields...).
		Render(&buffer)
	// write the Init method (static)
	buffer.WriteString(`

func (ct *ComponentsTable) Init(em *EntityManager) {
	ct.allocate()
	ct.linkEntityManager(em)
}

`)
	// write the allocate() function
	allocateStatements := make([]Code, len(componentNames))
	for i, componentName := range componentNames {
		allocateStatements[i] = Id("ct").Dot(componentName).
			Op("=").Op("&").Id(componentStructName(componentName)).Values()
	}
	Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("allocate").
		Params().
		Block(allocateStatements...).
		Render(&buffer)
	buffer.WriteString("\n\n")

	// write the linkEntityManager() function
	linkStatements := make([]Code, len(componentNames))
	for i, componentName := range componentNames {
		linkStatements[i] = Id("ct").Dot(componentName).Dot("em").
			Op("=").Id("em")
	}
	Func().
		Params(Id("ct").Op("*").Id("ComponentsTable")).
		Id("linkEntityManager").
		Params(Id("em").Op("*").Id("EntityManager")).
		Block(linkStatements...).
		Render(&buffer)
	return buffer.String()
}
