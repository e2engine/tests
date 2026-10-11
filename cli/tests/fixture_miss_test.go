package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestFixtureMiss(t *testing.T) {
	dataDir := "../data/fixture-miss"

	tests := []struct {
		name        string
		services    []string
		checkParams func(id string) *checkTestExecutionParams
	}{
		{
			name:     "fixture-miss-http-path",
			services: []string{"gateway-http-path"},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name:     "fixture-miss-http-body",
			services: []string{"gateway-http-body"},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name:     "fixture-miss-grpc-message",
			services: []string{"gateway-grpc-message"},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name:     "fixture-miss-grpc-method",
			services: []string{"gateway-grpc-method"},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name:     "fixture-miss-http-path-multiple",
			services: []string{"gateway-http-path"},
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				defaultTestTimeout,
			)
			defer cancel()

			h := harness.New(t, "../bin/cli", "", path.Join(t.TempDir(), "data"))
			serviceCtx, serviceCancel := context.WithCancel(context.Background())
			t.Cleanup(serviceCancel)
			for _, service := range test.services {
				s := h.StartService(serviceCtx, service)
				t.Cleanup(func() {
					s.Stop(t)
				})
			}

			// create environment
			environmentID := h.CreateEnvironment(
				ctx,
				test.name,
				path.Join(dataDir, test.name+".env.yml"),
			)

			// create test
			testID := h.CreateTest(
				ctx,
				test.name,
				path.Join(dataDir, test.name+".test.yml"),
			)

			// execute test
			t.Setenv("E2ENGINE_RUNNER_WORKER_REPORTER_DUMP_PAYLOAD", "true")
			testExecutionID := h.RunTest(
				ctx,
				testID,
				environmentID,
			)

			// get testexecution by ID
			actualTestExecution := h.GetTestExecution(ctx, testExecutionID)

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
}
