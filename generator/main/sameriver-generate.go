package main

import (
	"flag"
	"fmt"
	"github.com/dt-rush/sameriver/generator/generate"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const DEFAULT_ENGINE_DIR = "$GOPATH/src/github.com/dt-rush/sameriver/engine"
const DEFAULT_OUTPUT_DIR = "/tmp/sameriver"

var engineDir = flag.String("enginedir", DEFAULT_ENGINE_DIR,
	"the directory sameriver has been cloned to")

var gameDir = flag.String("gamedir", "",
	"the sameriver/ directory of the game we're generating engine source files for")

var outputDir = flag.String("outputdir", DEFAULT_OUTPUT_DIR,
	"the directory in which to output the generated files")

func main() {
	parseArgs()
	runGenerate()
}

func parseArgs() {
	flag.Parse()
	*engineDir = pathResolve(*engineDir)
	fmt.Printf("sameriver/engine/ dir is: %s\n", *engineDir)
	if *gameDir != "" {
		*gameDir = pathResolve(*gameDir)
		fmt.Printf("${yourgame}/../sameriver dir is: %s\n", *gameDir)
	}
	if *outputDir == DEFAULT_OUTPUT_DIR {
		os.MkdirAll(DEFAULT_OUTPUT_DIR, os.ModePerm)
	} else {
		*outputDir = pathResolve(*outputDir)
	}
	fmt.Printf("output dir is: %s\n", *outputDir)
}

func runGenerate() {
	g := generate.NewGenerateProcess(*engineDir, *gameDir, *outputDir)
	g.Run(generate.TargetsCollection{
		"events":     g.GenerateEventFiles,
		"components": g.GenerateComponentFiles,
		"world":      g.GenerateWorldFiles,
	})
	g.PrintReport()
	if g.HadErrors() {
		fmt.Println("Errors were encountered, exiting 1")
		os.Exit(1)
	} else {
		g.OutputFiles()
		if *gameDir != "" {
			g.CopyFiles()
		}
	}
	fmt.Printf("Done.\n")
}

func pathResolve(path string) string {
	resolved := path
	var err error
	// replace $GOPATH (from DEFAULT_ENGINE_DIR)
	resolved = strings.Replace(resolved, "$GOPATH", os.Getenv("GOPATH"), 1)
	// resolve replace tilde at start with home directory
	if rune(resolved[0]) == '~' {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		resolved = strings.Replace(resolved, "~", usr.HomeDir, 1)
	}
	// resolve relative paths
	if rune(resolved[0]) != '/' {
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			panic(err)
		}
	}
	return resolved
}
