package build

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

func (g *GenerateProcess) GenerateEventFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]string,
	moreTargets TargetsCollection) {

	// read the events.go file as an ast.File
	srcFileName := fmt.Sprintf("%s/events.go", g.engineDir)
	eventsAst, _, err := g.ReadSourceFile(srcFileName)
	if err != nil {
		msg := fmt.Sprintf("failed to generate ast.File for %s", srcFileName)
		return msg, err, nil, nil
	}
	// traverse the declarations in the ast.File to get the event names
	eventNames := getEventNames(srcFileName, eventsAst)
	if len(eventNames) == 0 {
		msg := fmt.Sprintf("no structs with name matching .*Event found in %s\n",
			srcFileName)
		return msg, nil, nil, nil
	}
	// generate source files
	sourceFiles = make(map[string]string)
	// generate enum source file
	sourceFiles["events_enum.go"] = generateEventsEnumFile(eventNames)
	// return
	return "generated", nil, sourceFiles, nil
}

func getEventNames(srcFile string, astFile *ast.File) []string {
	eventNames := make([]string, 0)
	for _, d := range astFile.Decls {
		decl, ok := d.(*(ast.GenDecl))
		if !ok {
			continue
		}
		if decl.Tok != token.TYPE {
			continue
		}
		name := decl.Specs[0].(*ast.TypeSpec).Name.Name
		if validName, _ := regexp.MatchString(".+Event", name); !validName {
			fmt.Printf("type %s in %s does not match regexp for an event "+
				"type (\".+Event\"). Will not include in generated files.\n",
				name, srcFile)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	return eventNames
}

func generateEventsEnumFile(eventNames []string) string {
	// for each event name, create an uppercase const name
	constNames := make(map[string]string)
	for _, eventName := range eventNames {
		eventNameStem := strings.Replace(eventName, "Event", "", 1)
		constNames[eventName] = strings.ToUpper(eventNameStem) + "_EVENT"
	}
	// generate the source file
	var buffer bytes.Buffer
	// write the top of the file
	buffer.WriteString("//\n//\n//\n// THIS FILE IS GENERATED\n//\n//\n//\n")
	buffer.WriteString("package engine\n\n")
	buffer.WriteString("type EventType int\n\n")
	buffer.WriteString(fmt.Sprintf(
		"const N_EVENT_TYPES = %d\n\n", len(eventNames)))
	// write the enum
	buffer.WriteString("const (\n")
	for _, eventName := range eventNames {
		buffer.WriteString(fmt.Sprintf(
			"\t%s = EventType(iota)\n",
			constNames[eventName]))
	}
	buffer.WriteString(")\n\n")
	// write the enum->string function
	buffer.WriteString("var EVENT_NAMES = map[EventType]string{\n")
	for _, eventName := range eventNames {
		buffer.WriteString(fmt.Sprintf(
			"\t%s: \"%s\",\n",
			constNames[eventName],
			constNames[eventName]))
	}
	buffer.WriteString("}")
	return buffer.String()
}
