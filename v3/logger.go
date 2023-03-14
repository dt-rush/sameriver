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
		log.Lmicroseconds)
	if os.Getenv("DISABLE_LOGGER") == "true" {
		logger.SetOutput(io.Discard)
	}
	return logger
}()
