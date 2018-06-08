package generate

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strings"
)

func (g *GenerateProcess) GenerateEventFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]*File) {

	// read the ${gameDir}/sameriver/events.go file as an ast.File
	srcFileName := path.Join(g.gameDir, "sameriver", "custom_events.go")
	eventsAst, _, err := readSourceFile(srcFileName)
	if err != nil {
		msg := fmt.Sprintf("failed to generate ast.File for %s", srcFileName)
		return msg, err, nil
	}
	// traverse the declarations in the ast.File to get the event names
	eventNames := getEventNames(srcFileName, eventsAst)
	sort.Strings(eventNames)
	if len(eventNames) == 0 {
		msg := fmt.Sprintf("no structs with name matching .*Event found in %s\n",
			srcFileName)
		return msg, nil, nil
	}
	// generate source files
	sourceFiles = make(map[string]*File)
	// generate enum source file
	sourceFiles["events_enum.go"] = generateEventsEnumFile(eventNames)
	// return
	return "generated", nil, sourceFiles
}

func getEventNames(srcFile string, astFile *ast.File) (
	eventNames []string) {
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
				name, srcFile)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	return eventNames
}

func generateEventsEnumFile(eventNames []string) *File {
	// for each event name, create an uppercase const name
	constNames := make(map[string]string)
	for _, eventName := range eventNames {
		eventNameStem := strings.Replace(eventName, "Event", "", 1)
		constNames[eventName] = strings.ToUpper(eventNameStem) + "_EVENT"
	}
	// generate the source file
	f := NewFile("engine")

	// typedef EventType int
	f.Type().Id("EventType").Int()

	// const N_EVENT_TYPES = ___
	f.Const().Id("N_EVENT_TYPES").Op("=").Lit(len(eventNames))

	// write the enum
	f.Const().DefsFunc(func(g *Group) {
		for _, eventName := range eventNames {
			g.Id(constNames[eventName]).Op("=").Iota()
		}
	})

	// write the enum->string function
	f.Var().Id("EVENT_NAMES").Op("=").
		Map(Id("EventType")).String().
		Values(DictFunc(func(d Dict) {
			for _, eventName := range eventNames {
				constName := constNames[eventName]
				d[Id(constName)] = Lit(constName)
			}
		}))

	return f

}