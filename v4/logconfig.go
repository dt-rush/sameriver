package sameriver

import (
	"fmt"
	"os"
	"time"

	"github.com/TwiN/go-color"

	"go.uber.org/atomic"
)

type PrintfLike func(format string, params ...any)

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

var logWarning = SubLogFunction(
	"WARNING", true,
	func(s string) string { return color.InYellow(color.InBold(s)) })

var logWarningRateLimited = func(ms int) PrintfLike {
	var flag atomic.Uint32
	return func(format string, params ...any) {
		if flag.CompareAndSwap(0, 1) {
			logWarning(format, params...)
			go func() {
				time.Sleep(time.Duration(ms) * time.Millisecond)
				flag.CompareAndSwap(1, 0)
			}()
		}
	}
}

var DEBUG_EVENTS = os.Getenv("DEBUG_EVENTS") == "true"
var logEvents = SubLogFunction(
	"Events", DEBUG_EVENTS,
	func(s string) string { return color.InWhiteOverPurple(s) })

var DEBUG_GOAP = os.Getenv("DEBUG_GOAP") == "true"
var logGOAPDebug = SubLogFunction(
	"GOAP", DEBUG_GOAP,
	func(s string) string { return s })

var DEBUG_RUNTIME_LIMITER = os.Getenv("DEBUG_RUNTIME_LIMITER") == "true"
var logRuntimeLimiter = SubLogFunction(
	"RuntimeLimiter", DEBUG_RUNTIME_LIMITER,
	func(s string) string { return s })
