package main

import (
	"flag"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var seconds = flag.Int("seconds", 10, "how many seconds to run for")
var printHash = flag.Bool("printhash", false,
	"boolean flag, supplied if you want to see what the spatial hash looks like")
