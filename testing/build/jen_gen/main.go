package main

import (
	"bytes"
	"flag"
	"fmt"
	// "github.com/dave/jennifer/jen"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	// "os"
	"os/user"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// helper function, thanks stack oveflow
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// flag which probably shouldn't need to provide
var engineDir = flag.String("enginedir",
	"",
	"the directory where the engine source code is")

// type definitions used by the generate process
type GenerateFunc func(
	target string, fset *token.FileSet) (
	message string, err error, moreTargets TargetsCollection)
type TargetsCollection map[string]GenerateFunc

// struct to hold data related to the generation
type GenerateProcess struct {
	fset             *token.FileSet
	sources          map[string]string
	messages         map[string]string
	errors           map[string]string
	rootTargets      TargetsCollection
	targetsProcessed []string
}

func NewGenerateProcess(rootTargets TargetsCollection) *GenerateProcess {
	g := GenerateProcess{}
	g.fset = token.NewFileSet()
	g.sources = make(map[string]string)
	g.messages = make(map[string]string)
	g.errors = make(map[string]string)
	g.rootTargets = rootTargets
	return &g
}

// used to run the targets and recursively run any more targets they produce
func (g *GenerateProcess) runTargets(targets TargetsCollection) {
	for target, f := range targets {
		fmt.Printf("----- running target: %s -----\n", target)
		message, err, moreTargets := f(target, g.fset)
		g.messages[target] = message
		g.errors[target] = fmt.Sprintf("%v", err)
		g.targetsProcessed = append(g.targetsProcessed, target)
		if moreTargets != nil {
			g.runTargets(moreTargets)
		}
	}
}

// used to run the process
func (g *GenerateProcess) Run() {
	g.runTargets(g.rootTargets)
}

// used to display a summary at the end
func (g *GenerateProcess) PrintReport() {
	fmt.Println("GENERATE PROCESS REPORT:\n===\n")
	for _, target := range g.targetsProcessed {
		fmt.Printf("## %s\n", target)
		msg := g.messages[target]
		fmt.Printf("message: %s\n", msg)
		err := g.errors[target]
		if err != "" {
			fmt.Printf("error: %s\n", err)
		}
		fmt.Println()
	}
}

//
// Generation functions
//

func generateComponentEnum(
	target string, fset *token.FileSet) (
	message string, err error, moreTargets TargetsCollection) {

	return "TODO", nil, nil
}

func generateComponentsTable(
	target string, fset *token.FileSet) (
	message string, err error, moreTargets TargetsCollection) {

	return "TODO", nil, nil
}

// generate the source files for each component (position_component.go, etc.)
func generateComponents(
	target string, fset *token.FileSet) (
	message string, err error, moreTargets TargetsCollection) {

	// TODO
	// get component set and for each component, create a target function
	// which will generate its source

	// TODO: return the targets
	return "TODO", nil, nil
}

//
// type CollisionEvent struct {
// 	EntityA EntityToken
// 	EntityB EntityToken
// }
//
// type SpawnRequest struct {
// 	EntityType int
// 	Position   [2]int16
// 	Active     bool
// }
//
// TO:
//
// type EventType int
//
// const N_EVENT_TYPES = 2
// const (
// 	EVENT_TYPE_COLLISION     = EventType(iota)
// 	EVENT_TYPE_SPAWN_REQUEST = EventType(iota)
// )
//
// var EVENT_NAMES = map[EventType]string{
// 	COLLISION_EVENT:     "collision event",
// 	SPAWN_REQUEST_EVENT: "spawn request event",
// }
func generateEventEnum(
	target string, fset *token.FileSet) (
	message string, err error, moreTargets TargetsCollection) {

	// read the events.go file and get it as an ast.File
	eventsDefFile := fmt.Sprintf("%s/events.go", *engineDir)
	var buffer bytes.Buffer
	eventsDefSrc, err := ioutil.ReadFile(eventsDefFile)
	if err != nil {
		msg := fmt.Sprintf("could not read %s", eventsDefFile)
		return msg, err, nil
	}
	astFile, err := parser.ParseFile(fset, "",
		eventsDefSrc, parser.AllErrors)
	if err != nil {
		msg := fmt.Sprintf("failed to generate astFile for %s",
			eventsDefFile)
		return msg, err, nil
	}
	// traverse the declarations in the ast.File to get the event names
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
				name, eventsDefFile)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	// generate the source file
	constNames := make(map[string]string)
	buffer.WriteString(`// THIS FILE IS GENERATED
package engine

type EventType int

`)
	buffer.WriteString(fmt.Sprintf(
		"const N_EVENT_TYPES = %d\n\n", len(eventNames)))
	buffer.WriteString("const (\n")
	for _, eventName := range eventNames {
		eventNameStem := strings.Replace(eventName, "Event", "", 1)
		constNames[eventName] = strings.ToUpper(eventNameStem) + "_EVENT"
		buffer.WriteString(fmt.Sprintf(
			"\t%s = EventType(iota)\n",
			constNames[eventName]))
	}
	buffer.WriteString(")\n\n")
	buffer.WriteString("var EVENT_NAMES = map[EventType]string{\n")
	for _, eventName := range eventNames {
		buffer.WriteString(fmt.Sprintf(
			"\t%s: \"%s\",\n",
			constNames[eventName],
			constNames[eventName]))
	}
	buffer.WriteString("}")
	fmt.Printf("================")
	fmt.Printf("GENERATED SOURCE")
	fmt.Printf("================\n\n")
	fmt.Println(buffer.String())
	fmt.Printf("================")
	fmt.Printf("================")
	fmt.Printf("================\n\n")
	return "generated", nil, nil
}

//
// main
//
func main() {
	flag.Parse()
	if *engineDir == "" {
		usr, _ := user.Current()
		homedir := usr.HomeDir
		*engineDir = fmt.Sprintf("%s/go/src/github.com/dt-rush/donkeys-qquest/engine", homedir)
	}
	fmt.Printf("Engine dir is: %s\n",
		*engineDir)

	var ROOT_TARGETS = TargetsCollection{
		"event_enum.go":     generateEventEnum,
		"component_enum.go": generateComponentEnum,
		"components":        generateComponents,
	}

	g := NewGenerateProcess(ROOT_TARGETS)
	g.Run()
	g.PrintReport()
}
