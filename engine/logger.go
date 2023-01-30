package engine

import (
	"fmt"
	"log"
	"os"
)

// global scope package singleton
var Logger = log.New(
	os.Stdout,
	"",
	log.Lmicroseconds)

// produce a function which will act like fmt.Sprintf but be silent or not
// based on a supplied boolean value (below the function definition in this
// file you can find all of them used)
func SubLogFunction(
	moduleName string, flag bool) func(s string, params ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", moduleName)
	return func(format string, params ...interface{}) {
		switch {
		case !flag:
			return
		case len(params) == 0:
			Logger.Printf(prefix + format)
		default:
			Logger.Printf(prefix+format, params...)
		}
	}
}
