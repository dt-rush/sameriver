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
	"",
	log.Lmicroseconds)
