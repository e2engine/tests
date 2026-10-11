package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestServiceCalledAtLeastOnce(t *testing.T) {
	dataDir := "../data/called-at-least-once"

	tests := []struct {
		name        string
		checkParams func(id string) *checkTestExecutionParams
	}{
		{
			name: "minimal-mocked-called-at-least-once",
			checkParams: func(id string) *checkTestExecutionParams {
				return &checkTestExecutionParams{
					useHTTPRequestHeaders: map[string]string{
						"E2engine-Test-Execution-Id": id,
					},
				}
			},
		},
		{
			name: "misspelled-fails-validation",
		},
		{
			name: "one-of-two-mocked-not-called",
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
