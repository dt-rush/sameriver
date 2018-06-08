package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"sort"
)

func (g *GenerateProcess) GenerateEventFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]*File) {

	// generate source files
	sourceFiles = make(map[string]*File)

	// seed file is the file in ${gameDir}/sameriver that we'll generate
	// engine code from
	seedFile := path.Join(g.gameDir, "custom_events.go")
	// engine base events file is the file in engineDir which holds the
	// base events which are integrated into the minimal reuqirements of
	// the engine
	engineBaseEventsFile := path.Join(g.engineDir, "base_events.go")

	// read from files
	var eventNames []string
	if g.gameDir != "" {
		eventNames = g.getEventNames(seedFile)
	}
	eventNames = append(eventNames, g.getEventNames(engineBaseEventsFile)...)
	sort.Strings(eventNames)

	// generate enum source file
	sourceFiles["events_enum.go"] = generateEventsEnumFile(eventNames)
	// return
	return "generated", nil, sourceFiles
}

func (g *GenerateProcess) getEventNames(srcFileName string) (
	eventNames []string) {

	astFile, _ := readSourceFile(srcFileName)

	// for each declaration in the source file
	for _, d := range astFile.Decls {
		// cast to generic declaration
		decl, ok := d.(*(ast.GenDecl))
		if !ok {
			continue
		}
		// if not a type declaration, continue
		if decl.Tok != token.TYPE {
			continue
		}
		// get the name of the type
		name := decl.Specs[0].(*ast.TypeSpec).Name.Name
		// if it's not a .+Event name, continue
		if validName, _ := regexp.MatchString(".+Event", name); !validName {
			fmt.Printf("type %s in %s does not match regexp for an event "+
				"type (\".+Event\"). Will not include in generated files.\n",
				name, srcFileName)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	return eventNames
}
