package config

import (
	"testing"
)

func TestMustLoad(t *testing.T) {
	var got *Config = MustLoad()
	if got == nil {
		t.Fatalf("MustLoad() error = failed to load config")
	}
}
