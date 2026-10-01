package util

import (
	"testing"

	"github.com/e2engine/core/model"
)

func AssertEqualTestSuiteExecutions(t *testing.T, expected, actual *model.TestSuiteExecution) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Fatalf("Expected nil, got %+v", actual)
	}
	if actual == nil {
		t.Fatalf("Expected %+v, got nil", expected)
	}

	if expected.ID != actual.ID {
		t.Fatalf("Expected ID %s, got %s", expected.ID, actual.ID)
	}
	if expected.TestSuiteID != actual.TestSuiteID {
		t.Fatalf("Expected TestSuiteID %s, got %s", expected.TestSuiteID, actual.TestSuiteID)
	}
	if expected.EnvironmentID != actual.EnvironmentID {
		t.Fatalf("Expected EnvironmentID %s, got %s", expected.EnvironmentID, actual.EnvironmentID)
	}
	if !expected.StartedAt.Equal(actual.StartedAt) {
		t.Fatalf("Expected StartedAt %v, got %v", expected.StartedAt, actual.StartedAt)
	}
	if !expected.FinishedAt.Equal(actual.FinishedAt) {
		t.Fatalf("Expected FinishedAt %v, got %v", expected.FinishedAt, actual.FinishedAt)
	}
	if expected.Status != actual.Status {
		t.Fatalf("Expected Status %s, got %s", expected.Status, actual.Status)
	}
	if expected.TestSuiteName != actual.TestSuiteName {
		t.Fatalf("Expected TestName %s, got %s", expected.TestSuiteName, actual.TestSuiteName)
	}
	if expected.EnvironmentName != actual.EnvironmentName {
		t.Fatalf("Expected EnvironmentName %s, got %s", expected.EnvironmentName, actual.EnvironmentName)
	}

	assertEqualTestSuiteExecutionSummary(t, expected.Summary, actual.Summary)
}

func assertEqualTestSuiteExecutionSummary(
	t *testing.T,
	expected,
	actual *model.TestSuiteExecutionSummary,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Fatalf("Expected Summary nil, got %+v", actual)
		return
	}
	if actual == nil {
		t.Fatalf("Expected Summary %+v, got nil", expected)
		return
	}

	if expected.Error != actual.Error {
		t.Fatalf("Expected Summary.Error %q, got %q", expected.Error, actual.Error)
	}
	if expected.Total != actual.Total {
		t.Fatalf("Expected Summary.Total %d, got %d", expected.Total, actual.Total)
	}
	if expected.Passed != actual.Passed {
		t.Fatalf("Expected Summary.Passed %d, got %d", expected.Passed, actual.Passed)
	}
	if expected.Failed != actual.Failed {
		t.Fatalf("Expected Summary.Failed %d, got %d", expected.Failed, actual.Failed)
	}
	if expected.Errors != actual.Errors {
		t.Fatalf("Expected Summary.Errors %d, got %d", expected.Errors, actual.Errors)
	}
}
