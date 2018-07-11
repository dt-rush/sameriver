package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"sort"
)

func (g *GenerateProcess) GenerateEventFiles(target string) GenerateOutput {

	output := GenerateOutput{
		generatedSourceFiles: make(map[string]GeneratedFile),
	}

	// seed file is the file in ${gameDir}/sameriver that we'll generate
	// engine code from
	customFile := path.Join(g.gameDir, "custom_events.go")
	// engine base events file is the file in engineDir which holds the
	// base events which are integrated into the minimal reuqirements of
	// the engine
	engineBaseEventsFile := path.Join(g.engineDir, "base_events.go")

	// read from files
	var eventNames []string
	// read event names from the engine base
	eventNames = append(eventNames, g.getEventNames(engineBaseEventsFile)...)
	// add in the custom event names from the game dir
	if g.gameDir != "" {
		customEventNames := g.getEventNames(customFile)
		// if there's a name collision, stop and return an error
		for _, baseEventName := range eventNames {
			for _, customEventName := range customEventNames {
				if baseEventName == customEventName {
					msg := fmt.Sprintf("event name collision between engine "+
						"and game custom code: %s appears twice\n",
						baseEventName)
					return GenerateOutput{msg, errors.New(msg), nil}
				}
			}
		}
		// if no error, we'll be here. append the custom names
		eventNames = append(eventNames, customEventNames...)
	}
	// sort the names
	sort.Strings(eventNames)

	// combine imports from seed file and engine base file
	var importStrings []string
	importStrings = append(importStrings,
		getImportStringsFromFile(engineBaseEventsFile)...)
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

	// generate enum source file
	output.generatedSourceFiles["events_enum.go"] = GeneratedFile{
		File:    generateEventsEnumFile(eventNames),
		Imports: importStrings,
	}
	// return
	return output
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
		if validName, _ := regexp.MatchString(".+Data", name); !validName {
			fmt.Printf("type %s in %s does not match regexp for an event "+
				"type (\".+Data\"). Will not include in generated files.\n",
				name, srcFileName)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	return eventNames
}
