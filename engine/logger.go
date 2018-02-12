/*
*
*
*
*/

package engine

import (
    "os"
    "log"
)

var TIME_LOGS = 0

var Logger = log.New (
    os.Stdout,
    "[DEBUG] ",
    log.Lmicroseconds & TIME_LOGS) 
