package tests

import (
	"context"
	"path"
	"sort"
	"testing"

	"github.com/e2engine/tests/cli/harness"
)

func TestSocketCRUDTestSuite(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTestTimeout,
	)
	defer cancel()

	t.Setenv("E2ENGINE_TRANSPORT_KIND", "socket")
	t.Setenv("E2ENGINE_TRANSPORT_ADDRESS", "127.0.0.1:9001")

	const (
		envSpecPath   = "../data/crud-testsuite/real-with-mocked-dependency.env.yml"
		test1SpecPath = "../data/crud-testsuite/gateway-calls-entitlements.test.yml"
		test2SpecPath = "../data/crud-testsuite/gateway-error.test.yml"
		tsSpecPath    = "../data/crud-testsuite/e2e.ts.yml"
		te1SpecPath   = "../data/crud-testsuite/http-real-mocked.te.yml"
		te2SpecPath   = "../data/crud-testsuite/gateway-error.te.yml"
		tseSpecPath   = "../data/crud-testsuite/e2e.tse.yml"
	)

	h := harness.New(t, "../bin/cli", "", path.Join(t.TempDir(), "data"))

	// create environment
	environmentID := h.CreateEnvironment(
		ctx,
		"real-with-mocked-dependency",
		envSpecPath,
	)

	// create test1
	test1ID := h.CreateTest(
		ctx,
		"gateway-calls-entitlements",
		test1SpecPath,
	)

	// create test2
	test2ID := h.CreateTest(
		ctx,
		"gateway-error",
		test2SpecPath,
	)

	// list testsuites, expect 0
	tss1 := h.ListTestSuites(ctx)
	if len(tss1) != 0 {
		t.Fatalf("Expected 0 testsuites, got: %d", len(tss1))
	}

	// create testsuite
	tsID := h.CreateTestSuite(
		ctx,
		"e2e",
		tsSpecPath,
	)

	// list testsuites, expect 1
	tss3 := h.ListTestSuites(ctx)
	if len(tss3) != 1 {
		t.Fatalf("Expected 1 testsuite, got: %d", len(tss3))
	}

	actualTS3 := tss3[0]
	checkTestSuite(t, tsID, tsSpecPath, &actualTS3)

	// get testsuite by ID
	actualTS4 := h.GetTestSuite(ctx, tsID)
	checkTestSuite(t, tsID, tsSpecPath, &actualTS4)

	// get testsuite by name
	actualTS5 := h.GetTestSuite(ctx, "e2e")
	checkTestSuite(t, tsID, tsSpecPath, &actualTS5)

	// list testsuite executions, expect 0
	tses1 := h.ListTestSuiteExecutions(ctx)
	if len(tses1) != 0 {
		t.Fatalf("Expected 0 test executions, got: %d", len(tses1))
	}

	// start gateway service
	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	t.Cleanup(serviceCancel)

	gateway := h.StartService(serviceCtx, "gateway-http-path")
	t.Cleanup(func() {
		gateway.Stop(t)
	})

	// execute testsuite
	tseID := h.RunTestSuite(
		ctx,
		tsID,
		environmentID,
	)

	// list testsuite executions, expect 1
	tses2 := h.ListTestSuiteExecutions(ctx)
	if len(tses2) != 1 {
		t.Fatalf("Expected 1 test execution, got: %d", len(tses2))
	}

	actualTSE2 := tses2[0]
	checkTestSuiteExecution(t, tseID, environmentID, tsID, tseSpecPath, &actualTSE2, true)

	// get testsuite execution by ID
	actualTSE4 := h.GetTestSuiteExecution(ctx, tseID)
	checkTestSuiteExecution(t, tseID, environmentID, tsID, tseSpecPath, &actualTSE4, false)

	// list test executions, expect 2
	tes5 := h.ListTestExecutions(ctx)
	if len(tes5) != 2 {
		t.Fatalf("Expected 2 test executions, got: %d", len(tes5))
	}

	// sort test executions by test name for comparison
	sort.Slice(tes5, func(i, j int) bool {
		return tes5[i].TestName < tes5[j].TestName
	})

	// compare test executions
	actualTE5 := tes5[0]
	checkTestExecution(t, actualTE5.ID, environmentID, test1ID, te1SpecPath, &actualTE5)

	actualTE6 := tes5[1]
	checkTestExecution(t, actualTE6.ID, environmentID, test2ID, te2SpecPath, &actualTE6)

	// delete testsuite execution by ID
	h.DeleteTestSuiteExecution(ctx, tseID, "")

	// list testsuite executions, expect 0
	tses7 := h.ListTestSuiteExecutions(ctx)
	if len(tses7) != 0 {
		t.Fatalf("Expected 0 test executions, got: %d", len(tses7))
	}
}
