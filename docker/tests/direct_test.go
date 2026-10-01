package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/docker/harness"
	"github.com/e2engine/tests/util"
)

func TestDirect(t *testing.T) {
	dataDir := "../data/direct"

	tests := []struct {
		name     string
		services []string
	}{
		{
			name: "direct-mocked-http",
		},
		{
			name: "direct-mocked-grpc",
		},
		{
			name: "direct-real-http",
			services: []string{
				"ok-http",
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

			h := harness.New(
				t,
				"",
				path.Join(t.TempDir(), "testdata"),
				"e2engine-network",
			)

			for _, service := range test.services {
				serviceCtx, serviceCancel := context.WithCancel(
					context.Background(),
				)
				t.Cleanup(serviceCancel)

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
					test.name+".env.yml",
				),
			)

			testID := h.CreateTest(
				ctx,
				test.name,
				path.Join(
					dataDir,
					test.name+".test.yml",
				),
			)

			testExecutionID := h.RunTest(
				ctx,
				testID,
				environmentID,
			)

			actualTestExecution := h.GetTestExecution(
				ctx,
				testExecutionID,
			)

			expectedTestExecution := util.Load[model.TestExecution](
				t,
				path.Join(
					dataDir,
					test.name+".yml",
				),
			)

			expectedTestExecution.ID = testExecutionID
			expectedTestExecution.TestID = testID
			expectedTestExecution.EnvironmentID = environmentID
			expectedTestExecution.StartedAt = actualTestExecution.StartedAt
			expectedTestExecution.FinishedAt = actualTestExecution.FinishedAt

			util.AssertEqualTestExecutions(
				t,
				&expectedTestExecution,
				&actualTestExecution,
			)
		})
	}
}
