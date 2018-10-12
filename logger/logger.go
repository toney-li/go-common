package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

type Logger struct {

}

func New(module string) *logrus.Logger {
	logger := logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.Formatter = &ApacheFormatter{}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.Out = os.Stdout

	// Only log the warning severity or above.
	logger.SetLevel(logrus.InfoLevel)
	logger.WithFields(logrus.Fields{FieldKeyModule: module})

	return logger
}
