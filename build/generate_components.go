package build

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func (g *GenerateProcess) GenerateComponentFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]string,
	moreTargets TargetsCollection) {

	// read the component_set.go file as an ast.File
	srcFileName := fmt.Sprintf("%s/component_set.go", g.engineDir)
	componentSetAst, componentSetSrc, err := g.ReadSourceFile(srcFileName)
	if err != nil {
		msg := fmt.Sprintf("failed to generate ast.File for %s", srcFileName)
		return msg, err, nil, nil
	}
	// get component set componentSetFields from src file
	componentSetFields, err := getComponentSetFields(
		componentSetSrc, componentSetAst)
	// generate source files
	sourceFiles = make(map[string]string)
	sourceFiles["components_enum.go"] =
		generateComponentsEnumFile(componentSetFields)
	sourceFiles["components_table.go"] =
		generateComponentsTableFile(componentSetFields)
	for componentName, componentType := range componentSetFields {
		filename := strings.ToLower(componentName) + "_component.go"
		sourceFiles[filename] = generateComponentFile(
			componentName, componentType)
	}
	// return
	return "generated", nil, sourceFiles, nil
}

func generateComponentsEnumFile(componentSetFields map[string]string) string {
	// for each component name, create an uppercase const name
	constNames := make(map[string]string)
	for componentName, _ := range componentSetFields {
		constNames[componentName] = strings.ToUpper(componentName) + "_COMPONENT"
	}
	// generate the source file
	var buffer bytes.Buffer
	// write the top of the file
	buffer.WriteString("//\n//\n//\n// THIS FILE IS GENERATED\n//\n//\n//\n")
	buffer.WriteString("package engine\n\n")
	buffer.WriteString("type ComponentType int\n\n")
	buffer.WriteString(fmt.Sprintf("const N_COMPONENT_TYPES = %d\n\n",
		len(componentSetFields)))
	// write the enum
	buffer.WriteString("const (\n")
	for componentName, _ := range componentSetFields {
		buffer.WriteString(fmt.Sprintf(
			"\t%s = componentType(iota)\n",
			constNames[componentName]))
	}
	buffer.WriteString(")\n\n")
	// write the enum->string function
	buffer.WriteString("var component_NAMES = map[componentType]string{\n")
	for componentName, _ := range componentSetFields {
		buffer.WriteString(fmt.Sprintf(
			"\t%s: \"%s\",\n",
			constNames[componentName],
			constNames[componentName]))
	}
	buffer.WriteString("}")
	return buffer.String()
}

func generateComponentsTableFile(
	componentSetFields map[string]string) string {

	return "TODO"
}

// generate the source files for each component (position_component.go, etc.)
func generateComponentFile(componentName string, componentType string) string {

	// TODO
	// get component set and for each component, create a target function
	// which will generate its source

	// TODO: return the targets
	return "TODO"
}

func getComponentSetFields(srcFile []byte, astFile *ast.File) (
	map[string]string, error) {

	for _, d := range astFile.Decls {
		decl, ok := d.(*(ast.GenDecl))
		if !ok {
			continue
		}
		if decl.Tok != token.TYPE {
			continue
		}
		typeSpec := decl.Specs[0].(*ast.TypeSpec)
		if typeSpec.Name.Name == "ComponentSet" {
			componentSetFields := make(map[string]string)
			for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {
				componentName := field.Names[0].Name
				componentType := string(
					srcFile[field.Type.Pos()-1 : field.Type.End()-1])
				fmt.Printf("found component: %s: %s\n",
					componentName, componentType)
				componentSetFields[componentName] = componentType
			}
			return componentSetFields, nil
		}
	}
	msg := fmt.Sprintf("no ComponentSet struct found in %s", srcFile)
	return nil, errors.New(msg)
}
