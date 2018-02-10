/*
*
*
*
*/

package utils

import (
    "fmt"
)

func DebugPrintln (msg string) {
    fmt.Println (msg)
}

func DebugPrintf (template string, args ...interface{}) {
    fmt.Printf (template, args...)
}
