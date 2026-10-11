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
		name         string
		services     []string
		httpHeader   bool
		grpcMetadata bool
	}{
		{
			name:       "direct-mocked-http",
			httpHeader: true,
		},
		{
			name:         "direct-mocked-grpc",
			grpcMetadata: true,
		},
		{
			name: "direct-real-http",
			services: []string{
				"ok-http",
			},
			httpHeader: true,
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
			if test.httpHeader {
				expectedTestExecution.Summary.Request.HTTP.Headers["E2engine-Test-Execution-Id"] =
					[]string{testExecutionID}
			}
			if test.grpcMetadata {
				expectedTestExecution.Summary.Request.GRPC.Metadata["e2engine-test-execution-id"] =
					[]string{testExecutionID}
			}

			util.AssertEqualTestExecutions(
				t,
				&expectedTestExecution,
				&actualTestExecution,
			)
		})
	}
}
