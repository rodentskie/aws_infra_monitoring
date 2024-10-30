package env

import (
	"testing"
)

func TestEnv(t *testing.T) {
	result := Env("works")
	if result != "Env works" {
		t.Error("Expected Env to append 'works'")
	}
}
