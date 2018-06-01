package build

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

func (g *GenerateProcess) GenerateComponentFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]string,
	moreTargets TargetsCollection) {

	// get needed info from src file
	componentSetFields, err := g.getComponentSetFields(
		fmt.Sprintf("%s/components/sameriver.go", g.gameDir),
		"ComponentSet")
	err = g.includeEngineBaseComponentSetFieldsInMap(componentSetFields)
	if err != nil {
		return "failed to include engine base component set", err, nil, nil
	}
	var componentNames []string
	for componentName, _ := range componentSetFields {
		componentNames = append(componentNames, componentName)
	}
	sort.Strings(componentNames)
	// generate source files
	sourceFiles = make(map[string]string)
	sourceFiles["components_enum.go"] =
		generateComponentsEnumFile(componentNames)
	sourceFiles["components_table.go"] =
		generateComponentsTableFile(componentNames)
	for componentName, componentType := range componentSetFields {
		filename := strings.ToLower(componentName) + "_component.go"
		sourceFiles[filename] = generateComponentFile(
			componentName, componentType)
	}
	// return
	return "generated", nil, sourceFiles, nil
}

func generateComponentsEnumFile(componentNames []string) string {
	// for each component name, create an uppercase const name
	constNames := make(map[string]string)
	for _, componentName := range componentNames {
		constNames[componentName] = strings.ToUpper(componentName) + "_COMPONENT"
	}
	// generate the source file
	var buffer bytes.Buffer

	Type().Id("ComponentType").Int().
		Render(&buffer)
	buffer.WriteString("\n\n")

	Const().Id("N_COMPONENT_TYPES").Op("=").Lit(len(componentNames)).
		Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum
	constDefs := make([]Code, len(componentNames))
	for i, componentName := range componentNames {
		constDefs[i] = Id(constNames[componentName]).Op("=").Iota()
	}
	Const().Defs(constDefs...).
		Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum->string function
	Var().Id("COMPONENT_NAMES").Op("=").
		Map(Id("ComponentType")).String().
		Values(DictFunc(func(d Dict) {
			for _, componentName := range componentNames {
				constName := constNames[componentName]
				d[Id(constName)] = Lit(constName)
			}
		})).
		Render(&buffer)
	return buffer.String()
}

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

func componentStructName(componentName string) string {
	return componentName + "Component"
}

// generate the source files for each component (position_component.go, etc.)
func generateComponentFile(componentName string, componentType string) string {

	// TODO
	// get component set and for each component, create a target function
	// which will generate its source

	// TODO: return the targets
	return "TODO"
}

func (g *GenerateProcess) getComponentSetFields(
	srcFileName string, structName string) (map[string]string, error) {

	// read the component_set.go file as an ast.File
	componentSetAst, componentSetSrcFile, err :=
		g.ReadSourceFile(srcFileName)
	if err != nil {
		fmt.Printf("failed to generate ast.File for %s",
			srcFileName)
		return nil, err
	}

	for _, d := range componentSetAst.Decls {
		decl, ok := d.(*(ast.GenDecl))
		if !ok {
			continue
		}
		if decl.Tok != token.TYPE {
			continue
		}
		typeSpec := decl.Specs[0].(*ast.TypeSpec)
		if typeSpec.Name.Name == structName {
			componentSetFields := make(map[string]string)
			for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {
				componentName := field.Names[0].Name
				if validName, _ :=
					regexp.MatchString(
						"[A-Z][a-z-A-Z]+", componentName); !validName {
					fmt.Printf("field %s in %s did not match regexp " +
						"[A-Z][a-z-A-Z]+ (exported field), and so won't " +
						"be considered a component")
					continue
				}
				componentType := string(
					componentSetSrcFile[field.Type.Pos()-1 : field.Type.End()-1])
				fmt.Printf("found component: %s: %s\n",
					componentName, componentType)
				componentSetFields[componentName] = componentType
			}
			return componentSetFields, nil
		}
	}
	msg := fmt.Sprintf("no ComponentSet struct found in %s",
		componentSetSrcFile)
	return nil, errors.New(msg)
}

func (g *GenerateProcess) includeEngineBaseComponentSetFieldsInMap(
	componentSetFields map[string]string) error {

	// read baes_component_set.go from the engine for merging
	engineBaseComponentSetFields, err := g.getComponentSetFields(
		fmt.Sprintf("%s/base_component_set.go", g.engineDir),
		"BaseComponentSet")
	if err != nil {
		return err
	}
	for componentName, componentType := range engineBaseComponentSetFields {
		componentSetFields[componentName] = componentType
	}
	return nil
}
