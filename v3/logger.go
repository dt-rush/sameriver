package sameriver

import (
	"log"
	"os"
)

// global scope package singleton
var Logger = log.New(
	os.Stdout,
	"",
	log.Lmicroseconds)
