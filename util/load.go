package util

import (
	"io"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func Load[T any](t *testing.T, path string) T {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var out T
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	return out
}
