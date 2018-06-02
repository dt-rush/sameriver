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

var engineDir = flag.String("enginedir",
	DEFAULT_ENGINE_DIR,
	"the directory where the sameriver/engine has been cloned to")
var gameDir = flag.String("gamedir",
	"",
	"the game/ directory of the game we're generating engine source files for")

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
	fmt.Printf("sameriver/engine/ dir is: %s\n", *engineDir)
	fmt.Printf("${yourgame}/game/ dir is: %s\n", *gameDir)

	g := build.NewGenerateProcess(*engineDir, *gameDir)
	g.Run(build.TargetsCollection{
		"events":     g.GenerateEventFiles,
		"components": g.GenerateComponentFiles,
	})
	g.PrintReport()
	g.PrintSourceFiles()
	if g.HadErrors() {
		os.Exit(1)
	}
}
