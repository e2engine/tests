package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestDirect(t *testing.T) {
	dataDir := "../data/direct"

	formats := []string{
		"yml",
		"json",
	}

	tests := []struct {
		name        string
		services    []string
		checkParams func(id string) *checkTestExecutionParams
	}{
		{
			name: "direct-mocked-http",
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name: "direct-mocked-grpc",
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useGRPCRequestMetadata: map[string]string{
						"e2engine-test-execution-id": id,
					},
				}
			},
		},
		{
			name: "direct-real-http",
			services: []string{
				"ok-http",
			},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name: "direct-real-grpc",
			services: []string{
				"ok-grpc",
			},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useGRPCRequestMetadata: map[string]string{
						"e2engine-test-execution-id": id,
					},
				}
			},
		},
	}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(
						context.Background(),
						defaultTestTimeout,
					)
					defer cancel()

					serviceCtx, serviceCancel := context.WithCancel(
						context.Background(),
					)
					t.Cleanup(serviceCancel)

					h := harness.New(
						t,
						"../bin/cli",
						"",
						path.Join(t.TempDir(), "data"),
					)

					for _, service := range test.services {
						s := h.StartService(serviceCtx, service)
						t.Cleanup(func() {
							s.Stop(t)
						})
					}

					environmentID := h.CreateEnvironment(
						ctx,
						test.name,
						path.Join(
							dataDir,
							test.name+".env."+format,
						),
					)

					testID := h.CreateTest(
						ctx,
						test.name,
						path.Join(
							dataDir,
							test.name+".test."+format,
						),
					)

					t.Setenv("E2ENGINE_RUNNER_WORKER_REPORTER_DUMP_PAYLOAD", "true")
					testExecutionID := h.RunTest(
						ctx,
						testID,
						environmentID,
					)

					actualTestExecution := h.GetTestExecution(
						ctx,
						testExecutionID,
					)

					if test.checkParams == nil {
						test.checkParams = func(id string) *checkTestExecutionParams {
							return nil
						}
					}

					checkTestExecution(
						t,
						testExecutionID,
						environmentID,
						testID,
						path.Join(dataDir, test.name+".yml"),
						&actualTestExecution,
						test.checkParams(testExecutionID),
					)
				})
			}
		})
	}
}
