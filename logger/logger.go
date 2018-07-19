package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func New(modul string) *log.Logger {
	logger := log.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.Formatter = &log.JSONFormatter{}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(log.InfoLevel)
	return logger
}
