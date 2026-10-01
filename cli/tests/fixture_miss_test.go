package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/cli/harness"
	"github.com/e2engine/tests/util"
)

func TestFixtureMiss(t *testing.T) {
	dataDir := "../data/fixture-miss"

	tests := []struct {
		name     string
		services []string
	}{
		{
			name:     "fixture-miss-http-path",
			services: []string{"gateway-http-path"},
		},
		{
			name:     "fixture-miss-http-body",
			services: []string{"gateway-http-body"},
		},
		{
			name:     "fixture-miss-grpc-message",
			services: []string{"gateway-grpc-message"},
		},
		{
			name:     "fixture-miss-grpc-method",
			services: []string{"gateway-grpc-method"},
		},
		{
			name:     "fixture-miss-http-path-multiple",
			services: []string{"gateway-http-path"},
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
			for _, service := range test.services {
				s := h.StartService(ctx, service)
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
