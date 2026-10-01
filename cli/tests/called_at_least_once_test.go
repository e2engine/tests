package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/cli/harness"
	"github.com/e2engine/tests/util"
)

func TestServiceCalledAtLeastOnce(t *testing.T) {
	dataDir := "../data/called-at-least-once"

	tests := []struct {
		name string
	}{
		{
			name: "minimal-mocked-called-at-least-once",
		},
		{
			name: "misspelled-called-at-least-once",
		},
		{
			name: "one-of-two-mocked-not-called",
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

			expectedTestExecution := util.Load[model.TestExecution](
				t,
				path.Join(dataDir, test.name+".yml"),
			)
			expectedTestExecution.ID = testExecutionID
			expectedTestExecution.TestID = testID
			expectedTestExecution.EnvironmentID = environmentID
			expectedTestExecution.StartedAt = actualTestExecution.StartedAt
			expectedTestExecution.FinishedAt = actualTestExecution.FinishedAt

			util.AssertEqualTestExecutions(t, &expectedTestExecution, &actualTestExecution)
		})
	}
}
