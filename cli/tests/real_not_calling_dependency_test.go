package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/cli/harness"
	"github.com/e2engine/tests/util"
)

func TestRealNotCallingDependency(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTestTimeout,
	)
	defer cancel()

	h := harness.New(t, "../bin/cli", "", path.Join(t.TempDir(), "data"))

	// create environment
	environmentID := h.CreateEnvironment(
		ctx,
		"real-not-calling-dependency",
		"../data/real-not-calling-dependency/real-not-calling-dependency.env.yml",
	)

	// create test
	testID := h.CreateTest(
		ctx,
		"real-not-calling-dependency",
		"../data/real-not-calling-dependency/real-not-calling-dependency.test.yml",
	)

	// start ok service
	okService := h.StartService(ctx, "ok-http")
	t.Cleanup(func() {
		okService.Stop(t)
	})

	// execute test, expect success
	testExecutionID := h.RunTest(
		ctx,
		testID,
		environmentID,
	)

	// get testexecution by ID
	actualTestExecution := h.GetTestExecution(ctx, testExecutionID)

	expectedTestExecution := util.Load[model.TestExecution](
		t,
		"../data/real-not-calling-dependency/real-not-calling-dependency.yml",
	)
	expectedTestExecution.ID = testExecutionID
	expectedTestExecution.TestID = testID
	expectedTestExecution.EnvironmentID = environmentID
	expectedTestExecution.StartedAt = actualTestExecution.StartedAt
	expectedTestExecution.FinishedAt = actualTestExecution.FinishedAt

	util.AssertEqualTestExecutions(t, &expectedTestExecution, &actualTestExecution)
}
