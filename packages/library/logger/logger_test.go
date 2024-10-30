package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	result := Logger("works")
	if result != "Logger works" {
		t.Error("Expected Logger to append 'works'")
	}
}
