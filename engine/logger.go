/*
*
*
*
 */

package engine

import (
	"log"
	"os"
)

var Logger = log.New(
	os.Stdout,
	"[DEBUG] ",
	log.Lmicroseconds)
