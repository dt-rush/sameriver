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

var Logger = log.New (
    os.Stdout,
    "[DEBUG] ",
    log.Lmicroseconds & DEBUG_TIMES)
