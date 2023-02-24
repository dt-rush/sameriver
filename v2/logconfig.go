package sameriver

import (
	"fmt"

	"github.com/TwiN/go-color"
)

const LOG_EVENTS = false

var logEvent = SubLogFunction(
	"Events", LOG_EVENTS, func(s string) string { return color.InWhiteOverPurple(s) })

var logWarning = SubLogFunction(
	"WARNING", true, func(s string) string { return color.InYellow(color.InBold(s)) })

// produce a function which will act like fmt.Sprintf but be silent or not
// based on a supplied boolean value (below the function definition in this
// file you can find all of them used)
func SubLogFunction(
	moduleName string,
	flag bool,
	wrapper func(s string) string) func(s string, params ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", moduleName)
	return func(format string, params ...interface{}) {
		switch {
		case !flag:
			return
		case len(params) == 0:
			Logger.Printf(wrapper(prefix + format))
		default:
			Logger.Printf(wrapper(fmt.Sprintf(prefix+format, params...)))
		}
	}
}
