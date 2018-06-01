package main

import (
	"flag"
	"fmt"
	"github.com/dt-rush/donkeys-qquest/build"
	"os/user"
)

var engineDir = flag.String("enginedir",
	"",
	"the directory where the engine source code is")

func main() {
	flag.Parse()
	if *engineDir == "" {
		usr, _ := user.Current()
		homedir := usr.HomeDir
		*engineDir = fmt.Sprintf("%s/go/src/github.com/dt-rush/donkeys-qquest/engine", homedir)
	}
	fmt.Printf("Engine dir is: %s\n",
		*engineDir)

	g := build.NewGenerateProcess(*engineDir)

	var ROOT_TARGETS = build.TargetsCollection{
		"events":     g.GenerateEventFiles,
		"components": g.GenerateComponentFiles,
	}

	g.Run(ROOT_TARGETS)
	g.PrintReport()
	g.PrintSourceFiles()
}
