package util

import (
	"testing"
)

func TestContainsKey(t *testing.T) {
	m:=make(map[string]interface{})
	ContainsKey(m,"123")
	t.Log()
}
