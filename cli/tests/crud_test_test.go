package tests

import (
	"context"
	"path"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestCRUDTest(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTestTimeout,
	)
	defer cancel()

	const (
		envSpecPath  = "../data/crud-test/real-with-mocked-dependency.env.yml"
		testSpecPath = "../data/crud-test/gateway-calls-entitlements.test.yml"
		execSpecPath = "../data/crud-test/http-real-mocked.yml"
	)

	h := harness.New(t, "../bin/cli", "", path.Join(t.TempDir(), "data"))

	// list environments, expect 0
	environments1 := h.ListEnvironments(ctx)
	if len(environments1) != 0 {
		t.Fatalf("Expected 0 environments, got: %d", len(environments1))
	}

	// create environment
	environmentID := h.CreateEnvironment(
		ctx,
		"real-with-mocked-dependency",
		envSpecPath,
	)

	// list environments, expect 1
	environments3 := h.ListEnvironments(ctx)
	if len(environments3) != 1 {
		t.Fatalf("Expected 1 environment, got: %d", len(environments3))
	}

	actual3 := environments3[0]
	checkEnvironment(t, environmentID, envSpecPath, &actual3)

	// get environment by ID
	actual4 := h.GetEnvironment(ctx, environmentID)
	checkEnvironment(t, environmentID, envSpecPath, &actual4)

	// get environment by name
	actual5 := h.GetEnvironment(ctx, "real-with-mocked-dependency")
	checkEnvironment(t, environmentID, envSpecPath, &actual5)

	// list tests, expect 0
	tests1 := h.ListTests(ctx)
	if len(tests1) != 0 {
		t.Fatalf("Expected 0 tests, got: %d", len(tests1))
	}

	// create test
	testID := h.CreateTest(
		ctx,
		"gateway-calls-entitlements",
		testSpecPath,
	)

	// list tests, expect 1
	tests3 := h.ListTests(ctx)
	if len(tests3) != 1 {
		t.Fatalf("Expected 1 test, got: %d", len(tests3))
	}

	actualTest3 := tests3[0]
	checkTest(t, testID, testSpecPath, &actualTest3)

	// get test by ID
	actualTest4 := h.GetTest(ctx, testID)
	checkTest(t, testID, testSpecPath, &actualTest4)

	// get test by name
	actualTest5 := h.GetTest(ctx, "gateway-calls-entitlements")
	checkTest(t, testID, testSpecPath, &actualTest5)

	// list test executions, expect 0
	executions1 := h.ListTestExecutions(ctx)
	if len(executions1) != 0 {
		t.Fatalf("Expected 0 test executions, got: %d", len(executions1))
	}

	// start gateway service
	serviceCtx, serviceCancel := context.WithCancel(
		context.Background(),
	)
	t.Cleanup(serviceCancel)

	gateway := h.StartService(serviceCtx, "gateway-http-path")
	t.Cleanup(func() {
		gateway.Stop(t)
	})

	// execute test, expect success
	testExecutionID := h.RunTest(
		ctx,
		testID,
		environmentID,
	)

	// list executions, expect 1
	executions2 := h.ListTestExecutions(ctx)
	if len(executions2) != 1 {
		t.Fatalf("Expected 1 test execution, got: %d", len(executions2))
	}

	actualExecution2 := executions2[0]
	checkTestExecution(t, testExecutionID, environmentID, testID, execSpecPath, &actualExecution2)

	// get testexecution by ID
	actualTestExecution4 := h.GetTestExecution(ctx, testExecutionID)
	checkTestExecution(t, testExecutionID, environmentID, testID, execSpecPath, &actualTestExecution4)

	// delete testexecution by ID
	h.DeleteTestExecution(ctx, testExecutionID, "")

	// delete environment by ID
	h.DeleteEnvironment(ctx, environmentID, "real-with-mocked-dependency")

	// list environments, expect 0
	environments7 := h.ListEnvironments(ctx)
	if len(environments7) != 0 {
		t.Fatalf("Expected 0 environments, got: %d", len(environments7))
	}

	// delete test by ID
	h.DeleteTest(ctx, testID, "gateway-calls-entitlements")

	// list tests, expect 0
	tests7 := h.ListTests(ctx)
	if len(tests7) != 0 {
		t.Fatalf("Expected 0 tests, got: %d", len(tests7))
	}
}
