package generate

// takes the engine's base_component_set.go and applies it on top of the
// game's components/sameriver_component_set.go ComponentSet struct, generating
// various component-related code in the engine

import (
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
)

type ComponentSpec struct {
	Name               string
	Type               string
	HasDeepCopyMethod  bool
	DeepCopyMethodFile string
}

func (g *GenerateProcess) GenerateComponentFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]*jen.File) {

	sourceFiles = make(map[string]*jen.File)

	// seed file is the file in ${gameDir}/sameriver that we'll generate engine
	// code from
	seedFile := path.Join(g.gameDir, "custom_component_set.go")
	// engine base component set file is the file in engineDir which holds the
	// base components which all entities can have according to the minimal
	// requirements of the engine
	engineBaseComponentSetFile := path.Join(g.engineDir, "base_component_set.go")

	// read from files
	var components []ComponentSpec
	if g.gameDir != "" {
		components = g.getComponentSpecs(seedFile, "CustomComponentSet")
	}
	components = append(components, g.getComponentSpecs(
		engineBaseComponentSetFile, "BaseComponentSet")...)
	sort.Slice(components, func(i int, j int) bool {
		return strings.Compare(components[i].Name, components[j].Name) == -1
	})

	// generate source files
	sourceFiles["components_enum.go"] =
		generateComponentsEnumFile(components)
	sourceFiles["component_set.go"] =
		generateComponentSetFile(components)
	sourceFiles["components_table.go"] =
		generateComponentsTableFile(components)
	sourceFiles["component_read_methods.go"] =
		generateComponentReadMethodsFile(components)
	// return
	return "generated", nil, sourceFiles
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
				hasDeepCopyMethod, deepCopyMethodFile :=
					g.getDeepCopyMethod(componentName)
				components = append(components,
					ComponentSpec{
						componentName,
						componentType,
						hasDeepCopyMethod,
						deepCopyMethodFile})
			}
			return components
		}
	}
	// if we're here, we didn't find struct ComponentSet in the file
	fmt.Printf("no struct named ComponentSet not found in %s",
		componentSetSrcFile)
	return components
}

func (g *GenerateProcess) getDeepCopyMethod(
	componentName string) (bool, string) {

	// attempt to find deep_copy_${component}.go either in enginedir or gamedir
	expectedFileName := fmt.Sprintf("deep_copy_%s.go",
		strings.ToLower(componentName))
	engineDirFiles, err := ioutil.ReadDir(g.engineDir)
	if err != nil {
		panic(err)
	}
	var gameDirFiles []os.FileInfo
	if g.gameDir != "" {
		gameDirFiles, err = ioutil.ReadDir(g.gameDir)
		if err != nil {
			panic(err)
		}
	}
	allFiles := make([]os.FileInfo, 0)
	allFiles = append(allFiles, engineDirFiles...)
	allFiles = append(allFiles, gameDirFiles...)
	for _, file := range allFiles {
		if file.Name() == expectedFileName {
			return true, file.Name()
		}
	}
	return false, ""
}
