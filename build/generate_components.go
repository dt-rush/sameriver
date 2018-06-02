package build

// takes the engine's base_component_set.go and applies it on top of the
// game's components/sameriver_component_set.go ComponentSet struct, generating
// various component-related code in the engine

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

func (g *GenerateProcess) GenerateComponentFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]string) {

	// seed file is the file in the gameDir that we'll generate engine
	// code from
	seedFile := fmt.Sprintf("%s/components/sameriver_component_set.go",
		g.gameDir)
	// engine base component set file is the file in engineDir which holds the
	// base components which all entities can have according to the minimal
	// requirements of the engine
	engineBaseComponentSetFile := fmt.Sprintf("%s/base_component_set.go",
		g.engineDir)

	// get needed info from src file
	componentSetFields, err := getComponentSetFields(seedFile, "ComponentSet")
	if err != nil {
		msg := fmt.Sprintf("failed to process %s", seedFile)
		return msg, err, nil
	}
	err = includeEngineBaseComponentSetFieldsInMap(
		engineBaseComponentSetFile, componentSetFields)
	if err != nil {
		msg := fmt.Sprintf("failed to process %s", engineBaseComponentSetFile)
		return msg, err, nil
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
	return "generated", nil, sourceFiles
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

func getComponentSetFields(
	srcFileName string, structName string) (map[string]string, error) {

	// read the component_set.go file as an ast.File
	componentSetAst, componentSetSrcFile, err :=
		readSourceFile(srcFileName)
	if err != nil {
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

func includeEngineBaseComponentSetFieldsInMap(
	engineBaseComponentSetFile string,
	componentSetFields map[string]string) error {

	// read baes_component_set.go from the engine for merging
	engineBaseComponentSetFields, err := getComponentSetFields(
		engineBaseComponentSetFile,
		"BaseComponentSet")
	if err != nil {
		return err
	}
	for componentName, componentType := range engineBaseComponentSetFields {
		componentSetFields[componentName] = componentType
	}
	return nil
}
