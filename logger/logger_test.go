package logger

import (
	"testing"
	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {

	logger := New("logger")
	logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}
