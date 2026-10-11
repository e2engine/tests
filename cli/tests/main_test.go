package tests

import (
	"sort"
	"testing"
	"time"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/util"
)

const defaultTestTimeout = 30 * time.Second

type Service struct {
	Name string
	Args []string
}

func checkEnvironment(t *testing.T, id, path string, actual *model.Environment) {
	t.Helper()

	expected := util.Load[model.Environment](t, path)
	expected.ID = id
	expected.CreatedAt = actual.CreatedAt
	expected.UpdatedAt = actual.UpdatedAt

	util.AssertEqualEnvironments(t, &expected, actual)
}

func checkTest(t *testing.T, id, path string, actual *model.Test) {
	t.Helper()

	expected := util.Load[model.Test](t, path)
	expected.ID = id
	expected.CreatedAt = actual.CreatedAt
	expected.UpdatedAt = actual.UpdatedAt

	util.AssertEqualTests(t, &expected, actual)
}

func checkTestSuite(t *testing.T, id, path string, actual *model.TestSuite) {
	t.Helper()

	expected := util.Load[model.TestSuite](t, path)
	expected.ID = id
	expected.CreatedAt = actual.CreatedAt
	expected.UpdatedAt = actual.UpdatedAt

	util.AssertEqualTestSuites(t, &expected, actual)
}

type checkTestExecutionParams struct {
	useHTTPRequestHeaders  map[string]string
	useGRPCRequestMetadata map[string]string
}

func checkTestExecution(
	t *testing.T,
	id, envID, testID, path string,
	actual *model.TestExecution,
	params *checkTestExecutionParams,
) {
	t.Helper()

	expected := util.Load[model.TestExecution](t, path)
	expected.ID = id
	expected.EnvironmentID = envID
	expected.TestID = testID
	expected.StartedAt = actual.StartedAt
	expected.FinishedAt = actual.FinishedAt

	if params != nil {
		if len(params.useHTTPRequestHeaders) > 0 {
			for k, v := range params.useHTTPRequestHeaders {
				expected.Summary.Request.HTTP.Headers[k] = []string{v}
			}
		}
		if len(params.useGRPCRequestMetadata) > 0 {
			for k, v := range params.useGRPCRequestMetadata {
				expected.Summary.Request.GRPC.Metadata[k] = []string{v}
			}
		}
	}

	util.AssertEqualTestExecutions(t, &expected, actual)
}

type checkTestSuiteExecutionParams struct {
	useHTTPRequestHeaders            map[string]string
	useHTTPRequestExecutionIDHeaders map[string]string // test name -> received header value
}

func checkTestSuiteExecution(
	t *testing.T,
	id, envID, tsID, path string,
	actual *model.TestSuiteExecution,
	skipTests bool,
	params *checkTestSuiteExecutionParams,
) {
	t.Helper()

	expected := util.Load[model.TestSuiteExecution](t, path)
	expected.ID = id
	expected.EnvironmentID = envID
	expected.TestSuiteID = tsID
	expected.StartedAt = actual.StartedAt
	expected.FinishedAt = actual.FinishedAt

	if params != nil {
		if len(params.useHTTPRequestHeaders) > 0 {
			for k, v := range params.useHTTPRequestHeaders {
				for _, test := range expected.Tests {
					test.Summary.Request.HTTP.Headers[k] = []string{v}
				}
			}
		}

		for i := range expected.Tests {
			if executionID, ok := params.useHTTPRequestExecutionIDHeaders[expected.Tests[i].TestName]; ok {
				expected.Tests[i].Summary.Request.HTTP.Headers["E2engine-Test-Execution-Id"] =
					[]string{executionID}
			}
		}
	}

	util.AssertEqualTestSuiteExecutions(t, &expected, actual)

	if skipTests {
		return
	}

	if len(expected.Tests) != len(actual.Tests) {
		t.Fatalf("Expected %d tests, got %d", len(expected.Tests), len(actual.Tests))
	}

	// sort test executions by test name for comparison
	sort.Slice(expected.Tests, func(i, j int) bool {
		return expected.Tests[i].TestName < expected.Tests[j].TestName
	})
	sort.Slice(actual.Tests, func(i, j int) bool {
		return actual.Tests[i].TestName < actual.Tests[j].TestName
	})

	for i := range expected.Tests {
		expected.Tests[i].ID = actual.Tests[i].ID
		expected.Tests[i].EnvironmentID = actual.Tests[i].EnvironmentID
		expected.Tests[i].TestID = actual.Tests[i].TestID
		expected.Tests[i].StartedAt = actual.Tests[i].StartedAt
		expected.Tests[i].FinishedAt = actual.Tests[i].FinishedAt
		util.AssertEqualTestExecutions(t, &expected.Tests[i], &actual.Tests[i])
	}
}
