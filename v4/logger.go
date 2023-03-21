package sameriver

import (
	"io"
	"log"
	"os"
)

var Logger = func() *log.Logger {
	logger := log.New(
		os.Stdout,
		"",
		// add log.Lshortfile if you wanna track down what file is logging something
		log.Lmicroseconds)
	if os.Getenv("DISABLE_LOGGER") == "true" {
		logger.SetOutput(io.Discard)
	}
	return logger
}()
