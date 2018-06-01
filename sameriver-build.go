package main

import (
	"flag"
	"fmt"
	"github.com/dt-rush/sameriver/build"
	"os"
)

var engineDir = flag.String("enginedir",
	"",
	"the directory where the sameriver/engine has been cloned to")
var gameDir = flag.String("gamedir",
	"",
	"the game/ directory of the game we're generating engine source files for")

func main() {
	flag.Parse()
	if *engineDir == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *gameDir == "" {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Printf("sameriver/engine/ dir is: %s\n", *engineDir)
	fmt.Printf("${yourgame}/game/ dir is: %s\n", *gameDir)

	g := build.NewGenerateProcess(*engineDir, *gameDir)

	var ROOT_TARGETS = build.TargetsCollection{
		"events":     g.GenerateEventFiles,
		"components": g.GenerateComponentFiles,
	}

	g.Run(ROOT_TARGETS)
	g.PrintReport()
	g.PrintSourceFiles()
}
