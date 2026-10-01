package tests

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestValidateEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		stdout      string
		stderr      string
		expectError bool
	}{
		{
			name:   "invalid-resource-kind",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: invalid resource kind",
			expectError: true,
		},
		{
			name:        "invalid-resource-id",
			stdout:      "",
			stderr:      "invalid environment spec, fields: ID",
			expectError: true,
		},
		{
			name:        "resource-name-too-short",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Name",
			expectError: true,
		},
		{
			name:        "resource-name-too-long",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Name",
			expectError: true,
		},
		{
			name:        "invalid-resource-version",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Version",
			expectError: true,
		},
		{
			name:        "resource-description-too-short",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Description",
			expectError: true,
		},
		{
			name:   "no-services",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: environment spec must contain at least one service",
			expectError: true,
		},
		{
			name:   "duplicate-service-id",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, environment.service.id: service, " +
				"validation: service ids must be unique",
			expectError: true,
		},
		{
			name:   "duplicate-service-address",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, environment.service.id: service-two, " +
				"environment.service.listen.address: 127.0.0.1:8081, " +
				"validation: service addresses must be unique",
			expectError: true,
		},
		{
			name:        "service-id-too-short",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].ID",
			expectError: true,
		},
		{
			name:        "service-id-too-long",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].ID",
			expectError: true,
		},
		{
			name:        "invalid-service-kind",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Kind",
			expectError: true,
		},
		{
			name:        "invalid-service-mode",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Mode",
			expectError: true,
		},
		{
			name:        "invalid-service-address",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Address",
			expectError: true,
		},
		{
			name:        "invalid-service-target",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].HTTPTarget",
			expectError: true,
		},
		{
			name:        "invalid-grpc-target",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].GRPCTarget",
			expectError: true,
		},
		{
			name:   "real-http-with-grpc-target",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: " +
				"GRPC target is not allowed for HTTP service",
			expectError: true,
		},
		{
			name:   "real-grpc-without-target",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: " +
				"target is required for real service",
			expectError: true,
		},
		{
			name:   "real-grpc-with-http-target",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: " +
				"HTTP target is not allowed for GRPC service",
			expectError: true,
		},
		{
			name:        "http-with-proto",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: environment, validation: proto spec not allowed for HTTP service",
			expectError: true,
		},
		{
			name:        "grpc-without-proto",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: environment, validation: proto spec is required for GRPC service",
			expectError: true,
		},
		{
			name:        "grpc-external-empty-service",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.External.Service",
			expectError: true,
		},
		{
			name:        "grpc-internal-empty-package",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.Internal.Package",
			expectError: true,
		},
		{
			name:        "grpc-internal-empty-service",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.Internal.Service",
			expectError: true,
		},
		{
			name:        "grpc-internal-empty-method-name",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.Internal.Methods[0].Name",
			expectError: true,
		},
		{
			name:        "grpc-internal-empty-field-name",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.Internal.Methods[0].Request.Fields[0].Name",
			expectError: true,
		},
		{
			name:        "grpc-internal-invalid-field-type",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Proto.Internal.Methods[0].Request.Fields[0].Type",
			expectError: true,
		},
		{
			name:        "mocked-without-fixtures",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: environment, validation: mocked service must have at least one fixture",
			expectError: true,
		},
		{
			name:        "mocked-with-target",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: environment, validation: mocked service must not have a target",
			expectError: true,
		},
		{
			name:        "real-http-without-target",
			stdout:      "",
			stderr:      "invalid spec, spec.kind: environment, validation: target is required for real service",
			expectError: true,
		},
		{
			name:   "real-http-with-fixtures",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: fixtures are not allowed for real service",
			expectError: true,
		},
		{
			name:   "http-fixture-missing-http-when",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: HTTP service fixture must contain HTTP when and then",
			expectError: true,
		},
		{
			name:   "http-fixture-missing-http-then",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: HTTP service fixture must contain HTTP when and then",
			expectError: true,
		},
		{
			name:   "grpc-fixture-on-http-service",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: gRPC fixture is not allowed for HTTP service",
			expectError: true,
		},
		{
			name:   "grpc-fixture-missing-grpc-when",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: gRPC service fixture must contain gRPC when and then",
			expectError: true,
		},
		{
			name:   "grpc-fixture-missing-grpc-then",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: gRPC service fixture must contain gRPC when and then",
			expectError: true,
		},
		{
			name:   "http-fixture-on-grpc-service",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, " +
				"validation: HTTP fixture is not allowed for gRPC service",
			expectError: true,
		},
		{
			name:        "http-fixture-invalid-method",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].When.HTTP.Method",
			expectError: true,
		},
		{
			name:        "http-fixture-empty-path",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].When.HTTP.Path",
			expectError: true,
		},
		{
			name:        "http-fixture-status-too-low",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].Then.HTTP.Status",
			expectError: true,
		},
		{
			name:        "http-fixture-status-too-high",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].Then.HTTP.Status",
			expectError: true,
		},
		{
			name:        "http-fixture-invalid-request-json",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].When.HTTP.Body",
			expectError: true,
		},
		{
			name:        "http-fixture-invalid-response-json",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].Then.HTTP.Body",
			expectError: true,
		},
		{
			name:   "grpc-proto-both-external-and-internal",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: exactly one of external or " +
				"internal proto spec must be provided for GRPC service",
			expectError: true,
		},
		{
			name:   "grpc-proto-neither-external-nor-internal",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: exactly one of external or " +
				"internal proto spec must be provided for GRPC service",
			expectError: true,
		},
		{
			name:   "grpc-external-neither-file-nor-buf-module",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: exactly one of file or " +
				"buf_module must be provided for external GRPC proto spec",
			expectError: true,
		},
		{
			name:   "grpc-external-file-and-buf-module",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: exactly one of file or " +
				"buf_module must be provided for external GRPC proto spec",
			expectError: true,
		},
		{
			name:   "real-grpc-with-internal-proto",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: internal proto spec " +
				"is not allowed for real GRPC service",
			expectError: true,
		},
		{
			name:   "grpc-internal-without-methods",
			stdout: "",
			stderr: "invalid spec, spec.kind: environment, validation: internal proto spec " +
				"must contain at least one method",
			expectError: true,
		},
		{
			name:        "grpc-fixture-empty-method",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].When.GRPC.Method",
			expectError: true,
		},
		{
			name:        "grpc-fixture-invalid-status",
			stdout:      "",
			stderr:      "invalid environment spec, fields: Spec.Services[0].Fixtures[0].Then.GRPC.Status",
			expectError: true,
		},
		{
			name:   "grpc-fixture-invalid-path",
			stdout: "",
			stderr: "the specified path is inaccessible, " +
				"file.path: ../data/validate-environment/grpc-fixture-invalid-path.yml, " +
				"cause: open ../data/validate-environment/grpc-fixture-invalid-path.yml: " +
				"no such file or directory",
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
				"environment",
				filepath.Join("../data/validate-environment/", tt.name+".yml"),
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
