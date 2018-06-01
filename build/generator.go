package build

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

// type definitions used by the generate process
type GenerateFunc func(target string) (
	message string,
	err error,
	sourceFiles map[string]string,
	moreTargets TargetsCollection)
type TargetsCollection map[string]GenerateFunc

// struct to hold data related to the generation
type GenerateProcess struct {
	engineDir        string
	sourceFiles      map[string]string
	messages         map[string]string
	errors           map[string]string
	rootTargets      TargetsCollection
	targetsProcessed []string
}

func NewGenerateProcess(
	engineDir string) *GenerateProcess {

	g := GenerateProcess{}
	g.engineDir = engineDir
	g.sourceFiles = make(map[string]string)
	g.messages = make(map[string]string)
	g.errors = make(map[string]string)
	return &g
}

// used to run the targets and recursively run any more targets they produce
func (g *GenerateProcess) runTargets(targets TargetsCollection) {
	for target, f := range targets {
		fmt.Printf("----- running target: %s -----\n", target)
		message, err, sourceFiles, moreTargets := f(target)
		g.messages[target] = message
		g.errors[target] = fmt.Sprintf("%v", err)
		for filename, contents := range sourceFiles {
			g.sourceFiles[filename] = contents
		}
		g.targetsProcessed = append(g.targetsProcessed, target)
		if moreTargets != nil {
			g.runTargets(moreTargets)
		}
	}
}

// used to run the process
func (g *GenerateProcess) Run(rootTargets TargetsCollection) {
	g.runTargets(rootTargets)
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

func (g *GenerateProcess) PrintSourceFiles() {
	fmt.Println("Source file output:")
	for filename, contents := range g.sourceFiles {
		fmt.Printf("---\n%s\n---\n\n", filename)
		fmt.Printf("%s\n\n\n", contents)
	}
}

func (g *GenerateProcess) ReadSourceFile(srcFileName string) (
	*ast.File, []byte, error) {

	src, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		return nil, []byte{}, err
	}
	astFile, err := parser.ParseFile(
		token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		return nil, []byte{}, err
	}
	return astFile, src, nil
}
