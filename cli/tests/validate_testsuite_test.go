package tests

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestValidateTestSuite(t *testing.T) {
	tests := []struct {
		name        string
		stdout      string
		stderr      string
		expectError bool
	}{
		{
			name:   "invalid-resource-kind",
			stdout: "",
			stderr: "invalid spec, spec.kind: testsuite, " +
				"validation: invalid resource kind",
			expectError: true,
		},
		{
			name:        "no-selectors",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: testsuite\ntest selectors spec must contain at least one of ids, names, or tags",
			expectError: true,
		},
		{
			name:        "empty-id",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.IDs[0]",
			expectError: true,
		},
		{
			name:        "id-too-long",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.IDs[0]",
			expectError: true,
		},
		{
			name:        "empty-name",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.Names[0]",
			expectError: true,
		},
		{
			name:        "name-too-long",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.Names[0]",
			expectError: true,
		},
		{
			name:        "empty-tag",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.Tags[0]",
			expectError: true,
		},
		{
			name:        "tag-too-long",
			stdout:      "",
			stderr:      "invalid testsuite spec, fields: Spec.Selectors.Tags[0]",
			expectError: true,
		},
		{
			name:        "id-min-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "id-max-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "name-min-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "name-max-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "tag-min-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "tag-max-length",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "ids-and-names",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "names-and-tags",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "all-selectors",
			stdout:      "Testsuite spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:   "duplicate-id",
			stdout: "",
			stderr: "invalid spec, spec.kind: testsuite\ntest selectors must be unique, " +
				"testsuite.selector: ids, value: id",
			expectError: true,
		},
		{
			name:   "duplicate-name",
			stdout: "",
			stderr: "invalid spec, spec.kind: testsuite\ntest selectors must be unique, " +
				"testsuite.selector: names, value: test",
			expectError: true,
		},
		{
			name:   "duplicate-tag",
			stdout: "",
			stderr: "invalid spec, spec.kind: testsuite\ntest selectors must be unique, " +
				"testsuite.selector: tags, value: tag",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTestTimeout)
			defer cancel()

			h := harness.New(t, "../bin/cli", "", "")

			stdout, stderr, err := h.RunCommand(
				ctx,
				"validate",
				"testsuite",
				filepath.Join("../data/validate-testsuite/", tt.name+".yml"),
			)

			if tt.expectError {
				if err == nil {
					t.Fatalf(
						"Expected an error but got none\nstdout: %s\nstderr: %s",
						stdout,
						stderr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf(
						"RunCommand failed: %v\nstdout: %s\nstderr: %s",
						err,
						stdout,
						stderr,
					)
				}
			}

			if strings.TrimSpace(stdout) != strings.TrimSpace(tt.stdout) {
				t.Errorf("Expected stdout: %q, got: %q", tt.stdout, stdout)
			}
			if strings.TrimSpace(stderr) != strings.TrimSpace(tt.stderr) {
				t.Errorf("Expected stderr: %q, got: %q", tt.stderr, stderr)
			}
		})
	}
}
