package build

import (
	"fmt"
)

// type definitions used by the generate process
type GenerateFunc func(target string) (
	message string,
	err error,
	sourceFiles map[string]string)
type TargetsCollection map[string]GenerateFunc

// struct to hold data related to the generation
type GenerateProcess struct {
	engineDir        string
	gameDir          string
	sourceFiles      map[string]string
	messages         map[string]string
	errors           map[string]string
	rootTargets      TargetsCollection
	targetsProcessed []string
}

func NewGenerateProcess(
	engineDir string, gameDir string) *GenerateProcess {

	g := GenerateProcess{}
	g.engineDir = engineDir
	g.gameDir = gameDir
	g.sourceFiles = make(map[string]string)
	g.messages = make(map[string]string)
	g.errors = make(map[string]string)
	return &g
}

// used to run the targets
func (g *GenerateProcess) Run(targets TargetsCollection) {
	for target, f := range targets {
		fmt.Printf("----- running target: %s -----\n", target)
		message, err, sourceFiles := f(target)
		g.messages[target] = message
		g.errors[target] = fmt.Sprintf("%v", err)
		for filename, contents := range sourceFiles {
			g.sourceFiles[filename] = contents
		}
		g.targetsProcessed = append(g.targetsProcessed, target)
	}
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
		fmt.Printf("//\n//\n//\n// THIS FILE IS GENERATED\n//\n//\n//\n")
		fmt.Printf("package engine\n\n")
		fmt.Printf("%s\n\n\n", contents)
	}
}

func (g *GenerateProcess) HadErrors() bool {
	return len(g.errors) > 0
}
