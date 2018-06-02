package main

import (
	"flag"
	"fmt"
	"github.com/dt-rush/sameriver/build"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_ENGINE_DIR = "$GOPATH/src/github.com/dt-rush/sameriver/engine"
const DEFAULT_OUTPUT_DIR = "/tmp/sameriver"

var engineDir = flag.String("enginedir",
	DEFAULT_ENGINE_DIR,
	"the directory where the sameriver/engine has been cloned to")
var gameDir = flag.String("gamedir",
	"",
	"the game/ directory of the game we're generating engine source files for")
var outputDir = flag.String("outputdir",
	DEFAULT_OUTPUT_DIR,
	"the directory in which to output the generated files")

func pathResolve(path string) string {
	resolved := path
	var err error
	// replace $GOPATH (from DEFAULT_ENGINE_DIR)
	resolved = strings.Replace(resolved, "$GOPATH", os.Getenv("GOPATH"), 1)
	// resolve relative paths
	if rune(resolved[0]) != '/' {
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			panic(err)
		}
	}
	return resolved
}

func main() {
	flag.Parse()
	*engineDir = pathResolve(*engineDir)
	if *gameDir == "" {
		flag.Usage()
		os.Exit(1)
	} else {
		*gameDir = pathResolve(*gameDir)
	}
	if *outputDir == DEFAULT_OUTPUT_DIR {
		os.MkdirAll(DEFAULT_OUTPUT_DIR, os.ModePerm)
	}
	fmt.Printf("sameriver/engine/ dir is: %s\n", *engineDir)
	fmt.Printf("${yourgame}/game/ dir is: %s\n", *gameDir)

	g := build.NewGenerateProcess(*engineDir, *gameDir, *outputDir)
	g.Run(build.TargetsCollection{
		"events":     g.GenerateEventFiles,
		"components": g.GenerateComponentFiles,
	})
	g.PrintReport()
	if g.HadErrors() {
		fmt.Println("Errors were encountered, exiting 1")
		os.Exit(1)
	} else {
		g.OutputFiles()
	}
}
