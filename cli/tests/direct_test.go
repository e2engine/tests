package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/cli/harness"
	"github.com/e2engine/tests/util"
)

func TestDirect(t *testing.T) {
	dataDir := "../data/direct"

	formats := []string{
		"yml",
		"json",
	}

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
		{
			name: "direct-real-grpc",
			services: []string{
				"ok-grpc",
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
							test.name+"."+format,
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
		})
	}
}
