package generate

// takes the engine's base_component_set.go and applies it on top of the
// game's components/sameriver_component_set.go ComponentSet struct, generating
// various component-related code in the engine

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strings"
)

type ComponentSpec struct {
	Name string
	Type string
}

func (g *GenerateProcess) GenerateComponentFiles(target string) GenerateOutput {

	output := GenerateOutput{
		generatedSourceFiles: make(map[string]GeneratedFile),
	}

	// seed file is the file in ${gameDir}/sameriver that we'll generate engine
	// code from
	customFile := path.Join(g.gameDir, "custom_component_set.go")
	// engine base component set file is the file in engineDir which holds the
	// base components which all entities can have according to the minimal
	// requirements of the engine
	engineBaseComponentSetFile := path.Join(g.engineDir, "base_component_set.go")

	// read from files
	var components []ComponentSpec
	// read components from the engine base
	components = append(components, g.getComponentSpecs(
		engineBaseComponentSetFile, "BaseComponentSet")...)
	// add in the custom components from the game dir
	if g.gameDir != "" {
		customComponents := g.getComponentSpecs(customFile, "CustomComponentSet")
		// if there's a name collision, stop and return error
		for _, baseComponent := range components {
			for _, customComponent := range customComponents {
				if baseComponent.Name == customComponent.Name {
					msg := fmt.Sprintf("component name collision between "+
						"engine and game custom code: %s appears twice\n",
						baseComponent.Name)
					return GenerateOutput{msg, errors.New(msg), nil}
				}
			}
		}
		components = append(components, customComponents...)
	}
	// sort the names
	sort.Slice(components, func(i int, j int) bool {
		return strings.Compare(components[i].Name, components[j].Name) == -1
	})

	// combine imports from seed file and engine base file
	var importStrings []string
	importStrings = append(importStrings,
		getImportStringsFromFile(engineBaseComponentSetFile)...)
	if g.gameDir != "" {
		// if the import is already in the engine base file's imports,
		// skip inclusion, else add to the list of import strings
		for _, customImport := range getImportStringsFromFile(customFile) {
			inBaseImports := false
			for _, baseImport := range importStrings {
				if baseImport == customImport {
					inBaseImports = true
					break
				}
			}
			if !inBaseImports {
				importStrings = append(importStrings, customImport)
			}
		}
	}

	// generate source files
	output.generatedSourceFiles["components_enum.go"] = GeneratedFile{
		File:    generateComponentsEnumFile(components),
		Imports: make([]string, 0),
	}
	output.generatedSourceFiles["component_set.go"] = GeneratedFile{
		File:    generateComponentSetFile(components),
		Imports: importStrings,
	}
	output.generatedSourceFiles["components_table.go"] = GeneratedFile{
		File:    generateComponentsTableFile(components),
		Imports: importStrings,
	}
	output.generatedSourceFiles["entity_components_get_methods.go"] = GeneratedFile{
		File:    generateEntityComponentsGetMethodsFile(components),
		Imports: importStrings,
	}
	// return
	output.message = "generated"
	return output
}

func componentStructName(componentName string) string {
	return componentName + "Component"
}

func (g *GenerateProcess) getComponentSpecs(
	srcFileName string, structName string) (components []ComponentSpec) {

	// read the component_set.go file as an ast.File
	componentSetAst, componentSetSrcFile := readSourceFile(srcFileName)

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
			for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {
				componentName := field.Names[0].Name
				if validName, _ :=
					regexp.MatchString(
						"[A-Z][a-z-A-Z]+", componentName); !validName {
					fmt.Printf("field %s in %s did not match regexp "+
						"[A-Z][a-z-A-Z]+ (exported field), and so won't "+
						"be considered a component",
						componentName, srcFileName)
					continue
				}
				componentType := string(
					componentSetSrcFile[field.Type.Pos()-1 : field.Type.End()-1])
				if validType, _ :=
					regexp.MatchString(
						"\\\\*.+", componentName); !validType {
					fmt.Printf("%s's field type %s in %s is not pointer. "+
						"All ComponentSet members must be pointer type.\n",
						componentName, componentType, srcFileName)
					continue
				}
				componentType = strings.Replace(componentType, "*", "", 1)
				fmt.Printf("found component: %s: %s\n",
					componentName, componentType)
				components = append(components,
					ComponentSpec{
						componentName,
						componentType})
			}
			return components
		}
	}
	// if we're here, we didn't find struct ComponentSet in the file
	fmt.Printf("no struct named ComponentSet not found in %s",
		componentSetSrcFile)
	return components
}
