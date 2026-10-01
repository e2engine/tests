package tests

import (
	"context"
	"testing"

	coremodel "github.com/e2engine/core/model"
	"gopkg.in/yaml.v3"

	"github.com/e2engine/tests/cli/harness"
)

func TestVersion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTestTimeout)
	defer cancel()

	h := harness.New(t, "../bin/cli", "", "")

	stdout, stderr, err := h.RunCommand(ctx, "version")
	if err != nil {
		t.Fatalf("RunCommand failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	var versionOutput coremodel.Version
	if err := yaml.Unmarshal([]byte(stdout), &versionOutput); err != nil {
		t.Fatalf("Failed to unmarshal version output: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	if versionOutput.Version == "" {
		t.Errorf("Version is empty")
	}
	if versionOutput.GitCommit == "" {
		t.Errorf("GitCommit is empty")
	}
	if versionOutput.BuildTime.IsZero() {
		t.Errorf("BuildTime is zero")
	}
	if versionOutput.GoVersion == "" {
		t.Errorf("GoVersion is empty")
	}
	if versionOutput.Compiler == "" {
		t.Errorf("Compiler is empty")
	}
	if versionOutput.Platform == "" {
		t.Errorf("Platform is empty")
	}

	if stderr != "" {
		t.Errorf("Expected no stderr output, got: %s", stderr)
	}
}
