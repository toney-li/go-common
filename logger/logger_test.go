package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {

	logger := New("logger")
	logger.Info("A group of walrus emerges from the ocean")
}
