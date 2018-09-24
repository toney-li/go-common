package util

import (
	"testing"
	"fmt"
)

func TestRandomNumber(t *testing.T) {
	for {
		i := RandomNumber(7, 10, 1)
		fmt.Println(i)
	}
}
