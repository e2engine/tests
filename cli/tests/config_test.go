package tests

import (
	"context"
	"path"
	"strings"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		env            map[string]string
		expectedError  bool
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "default config",
			expectedStdout: "[]",
		},
		{
			name: "invalid transport kind",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND": "invalid_transport_kind",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"- Field \"Transport.Kind\": rule \"oneof\": " +
				"rule constraint violated (allowed=direct,socket)",
		},
		{
			name: "direct transport kind",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND": "direct",
			},
			expectedStdout: "[]",
		},
		{
			name: "direct transport kind with address",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "direct",
				"E2ENGINE_TRANSPORT_ADDRESS": "localhost:1234",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, " +
				"cause: invalid transport configuration, " +
				"validation: transport address is not allowed for direct transport",
		},
		{
			name: "socket transport kind with address",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "localhost:1234",
			},
			expectedStdout: "[]",
		},
		{
			name: "socket transport kind without address",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND": "socket",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, " +
				"cause: invalid transport configuration, " +
				"validation: transport address must be provided for socket transport",
		},
		{
			name: "socket transport kind with invalid address",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "invalid_address",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, " +
				"cause: cannot load config, " +
				"cause: cannot initialize configuration object, " +
				"cause: - Field \"Transport.Address\": rule \"networkaddress\": " +
				"invalid spec, validation: network address must be in host:port form, " +
				"cause: address invalid_address: missing port in address",
		},
		{
			name: "workers pool size",
			env: map[string]string{
				"E2ENGINE_RUNNER_WORKERS_POOL_SIZE": "4",
			},
			expectedStdout: "[]",
		},
		{
			name: "invalid workers pool size",
			env: map[string]string{
				"E2ENGINE_RUNNER_WORKERS_POOL_SIZE": "0",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"- Field \"Runner.WorkersPoolSize\": rule \"min\": " +
				"rule constraint violated (value=1)",
		},
		{
			name: "list max size min value",
			env: map[string]string{
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE": "1",
			},
			expectedStdout: "[]",
		},
		{
			name: "list max size invalid value",
			env: map[string]string{
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE": "0",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"- Field \"Output.ListMaxSize\": rule \"min\": " +
				"rule constraint violated (value=1)",
		},
		{
			name: "list max size max value",
			env: map[string]string{
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE": "50",
			},
			expectedStdout: "[]",
		},
		{
			name: "list max size above max value",
			env: map[string]string{
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE": "51",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"- Field \"Output.ListMaxSize\": rule \"max\": " +
				"rule constraint violated (value=50)",
		},
		{
			name: "invalid list max size type",
			env: map[string]string{
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE": "invalid",
			},
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"cannot set default value, field.name: Output.ListMaxSize, " +
				"env: E2ENGINE_OUTPUT_LIST_MAX_SIZE, cause: " +
				"cannot parse int, value: invalid, cause: " +
				"strconv.ParseInt: parsing \"invalid\": invalid syntax",
		},
		{
			name:           "invalid file config",
			configPath:     "../config/invalid_output_list_max_size.yml",
			expectedError:  true,
			expectedStdout: "",
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"- Field \"Output.ListMaxSize\": rule \"max\": " +
				"rule constraint violated (value=50)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				defaultTestTimeout,
			)
			defer cancel()

			h := harness.New(t, "../bin/cli", test.configPath, path.Join(t.TempDir(), "data"))

			for key, value := range test.env {
				t.Setenv(key, value)
			}

			stdout, stderr, err := h.RunCommand(ctx, "get", "environments", "-o", "yaml")
			if test.expectedError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
				}
			}

			if strings.TrimSpace(stdout) != test.expectedStdout {
				t.Fatalf("Expected stdout: %q, got: %q", test.expectedStdout, stdout)
			}
			if strings.TrimSpace(stderr) != test.expectedStderr {
				t.Fatalf("Expected stderr: %q, got: %q", test.expectedStderr, stderr)
			}
		})
	}
}

func TestConfig_ExecutionCheckTimeout(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectedError  bool
		expectedStderr string
	}{
		{
			name:  "valid execution check timeout",
			value: "10s",
		},
		{
			name:          "invalid execution check timeout",
			value:         "invalid",
			expectedError: true,
			expectedStderr: "cannot initialize service, cause: cannot load config, cause: " +
				"cannot initialize configuration object, cause: " +
				"cannot set default value, field.name: Runtime.ExecutionCheckTimeout, " +
				"env: E2ENGINE_RUNTIME_EXECUTION_CHECK_TIMEOUT, cause: " +
				"cannot parse duration, value: invalid, cause: " +
				"time: invalid duration \"invalid\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				defaultTestTimeout,
			)
			defer cancel()

			t.Setenv(
				"E2ENGINE_RUNTIME_EXECUTION_CHECK_TIMEOUT",
				test.value,
			)

			h := harness.New(
				t,
				"../bin/cli",
				"",
				path.Join(t.TempDir(), "data"),
			)

			stdout, stderr, err := h.RunCommand(
				ctx,
				"get",
				"environments",
				"-o",
				"yaml",
			)

			if test.expectedError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf(
					"Unexpected error: %v\nstdout: %s\nstderr: %s",
					err,
					stdout,
					stderr,
				)
			}

			if test.expectedError {
				if strings.TrimSpace(stdout) != "" {
					t.Fatalf(
						"Expected stdout: %q, got: %q",
						"",
						stdout,
					)
				}
			} else {
				if strings.TrimSpace(stdout) != "[]" {
					t.Fatalf(
						"Expected stdout: %q, got: %q",
						"[]",
						stdout,
					)
				}
			}

			if strings.TrimSpace(stderr) != test.expectedStderr {
				t.Fatalf(
					"Expected stderr: %q, got: %q",
					test.expectedStderr,
					stderr,
				)
			}
		})
	}
}

// TODO: add test for actual timeout resolution, simulate using real service response delay.
