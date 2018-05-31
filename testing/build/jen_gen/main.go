package main

import (
	"bytes"
	"flag"
	"fmt"
	// "github.com/dave/jennifer/jen"
	// "go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	// "os"
	"os/user"
	"reflect"
	"runtime"
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

	eventsDefFile := fmt.Sprintf("%s/events.go", *engineDir)
	var buffer bytes.Buffer
	eventsDefSrc, err := ioutil.ReadFile(eventsDefFile)
	if err != nil {
		msg := fmt.Sprintf("could not read %s", eventsDefFile)
		return msg, err, nil
	}
	astFile, err := parser.ParseFile(fset, "",
		eventsDefSrc, parser.AllErrors)
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
	fmt.Printf("Generating code from engine dir: %s ...\n",
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
