package tests

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestValidateTest(t *testing.T) {
	tests := []struct {
		name        string
		stdout      string
		stderr      string
		expectError bool
	}{
		{
			name:   "invalid-resource-kind",
			stdout: "",
			stderr: "invalid spec, spec.kind: test, " +
				"validation: invalid resource kind",
			expectError: true,
		},
		{
			name:        "request-missing",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nrequest must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "request-http-and-grpc",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nrequest must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "expect-missing",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nexpect must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "expect-http-and-grpc",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nexpect must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "request-expect-protocol-mismatch",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nrequest and expect protocols must match",
			expectError: true,
		},
		{
			name:        "empty-tag",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Tags[0]",
			expectError: true,
		},
		{
			name:        "tag-too-long",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Tags[0]",
			expectError: true,
		},
		{
			name:        "http-request-invalid-method",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.HTTP.Method",
			expectError: true,
		},
		{
			name:        "http-request-empty-url",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.HTTP.URL",
			expectError: true,
		},
		{
			name:        "http-request-invalid-url",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\ninvalid HTTP request URL, cause: parse \"http://[::1\": missing ']' in host",
			expectError: true,
		},
		{
			name:        "http-request-invalid-scheme",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nHTTP request URL scheme must be http or https",
			expectError: true,
		},
		{
			name:        "http-request-no-host",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\nHTTP request URL must contain a host",
			expectError: true,
		},
		{
			name:        "http-request-invalid-json",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.HTTP.Body",
			expectError: true,
		},
		{
			name:        "http-expect-status-too-low",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.HTTP.Status",
			expectError: true,
		},
		{
			name:        "http-expect-status-too-high",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.HTTP.Status",
			expectError: true,
		},
		{
			name:        "http-expect-invalid-json",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.HTTP.Body",
			expectError: true,
		},
		{
			name:        "grpc-request-invalid-target",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.GRPC.Target",
			expectError: true,
		},
		{
			name:        "grpc-request-empty-service",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.GRPC.Service",
			expectError: true,
		},
		{
			name:        "grpc-request-empty-method",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Request.GRPC.Method",
			expectError: true,
		},
		{
			name:        "grpc-expect-empty-status",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.GRPC.Status",
			expectError: true,
		},
		{
			name:        "call-no-protocol",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\ncall expectation must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "call-http-and-grpc",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: test\ncall expectation must contain exactly one of http or grpc",
			expectError: true,
		},
		{
			name:        "call-service-id-too-short",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].ServiceID",
			expectError: true,
		},
		{
			name:        "call-service-id-too-long",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].ServiceID",
			expectError: true,
		},
		{
			name:        "call-negative-count",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].Count",
			expectError: true,
		},
		{
			name:        "call-http-invalid-method",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].HTTP.Method",
			expectError: true,
		},
		{
			name:        "call-http-empty-path",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].HTTP.Path",
			expectError: true,
		},
		{
			name:        "call-http-invalid-json",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].HTTP.Body",
			expectError: true,
		},
		{
			name:        "call-grpc-empty-service",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].GRPC.Service",
			expectError: true,
		},
		{
			name:        "call-grpc-empty-method",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.Calls[0].GRPC.Method",
			expectError: true,
		},
		{
			name:        "call-count-zero",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-expect-status-min",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-expect-status-max",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "service-id-min-length",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "service-id-max-length",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name: "duplicate-tag",
			stderr: "invalid spec, spec.kind: test\ntest tags must be unique, " +
				"test.tag: tag",
			expectError: true,
		},
		{
			name:        "tag-min-length",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "tag-max-length",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-url-http",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-url-https",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-request-valid-json-scalar",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-request-all-optional-fields-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-call-all-optional-fields-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-request-all-optional-fields-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "http-request-valid-json-array",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-expect-invalid-status",
			stdout:      "",
			stderr:      "invalid test spec, fields: Spec.Expect.GRPC.Status",
			expectError: true,
		},
		{
			name:        "http-expect-valid-json-null",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-minimal-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-call-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-request-message-omitted-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-call-message-omitted-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "call-count-omitted-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-expect-message-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
		},
		{
			name:        "grpc-call-all-optional-fields-valid",
			stdout:      "Test spec is valid",
			stderr:      "",
			expectError: false,
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
				"test",
				filepath.Join("../data/validate-test/", tt.name+".yml"),
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
