package util

import (
	"strconv"
	"testing"

	"github.com/e2engine/core/model"
)

func AssertEqualTestExecutions(t *testing.T, expected, actual *model.TestExecution) {
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
		t.Errorf("Expected ID %s, got %s", expected.ID, actual.ID)
	}
	if expected.TestID != actual.TestID {
		t.Errorf("Expected TestID %s, got %s", expected.TestID, actual.TestID)
	}
	if expected.EnvironmentID != actual.EnvironmentID {
		t.Errorf("Expected EnvironmentID %s, got %s", expected.EnvironmentID, actual.EnvironmentID)
	}
	if !expected.StartedAt.Equal(actual.StartedAt) {
		t.Errorf("Expected StartedAt %v, got %v", expected.StartedAt, actual.StartedAt)
	}
	if !expected.FinishedAt.Equal(actual.FinishedAt) {
		t.Errorf("Expected FinishedAt %v, got %v", expected.FinishedAt, actual.FinishedAt)
	}
	if expected.Status != actual.Status {
		t.Errorf("Expected Status %s, got %s", expected.Status, actual.Status)
	}
	if expected.TestName != actual.TestName {
		t.Errorf("Expected TestName %s, got %s", expected.TestName, actual.TestName)
	}
	if expected.EnvironmentName != actual.EnvironmentName {
		t.Errorf("Expected EnvironmentName %s, got %s", expected.EnvironmentName, actual.EnvironmentName)
	}

	assertEqualTestExecutionSummary(t, expected.Summary, actual.Summary)
}

func assertEqualTestExecutionSummary(
	t *testing.T,
	expected,
	actual *model.TestExecutionSummary,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected Summary nil, got %+v", actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected Summary %+v, got nil", expected)
		return
	}

	if expected.Error != actual.Error {
		t.Errorf("Expected Summary.Error %q, got %q", expected.Error, actual.Error)
	}

	assertEqualTestExecutionRequestSummary(
		t,
		"Summary.Request",
		expected.Request,
		actual.Request,
	)

	assertEqualTestExecutionResponseSummary(
		t,
		"Summary.Response",
		expected.Response,
		actual.Response,
	)

	assertEqualAllServiceObservedCalls(
		t,
		"Summary.Calls",
		expected.Calls,
		actual.Calls,
	)

	assertEqualTestExecutionDeviations(
		t,
		"Summary.Deviations",
		expected.Deviations,
		actual.Deviations,
	)
}

func assertEqualTestExecutionRequestSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionRequestSummary,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected %s nil, got %+v", path, actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected %s %+v, got nil", path, expected)
		return
	}

	switch {
	case expected.HTTP != nil && actual.HTTP != nil:
		assertEqualHTTPRequest(
			t,
			path+".HTTP",
			expected.HTTP,
			actual.HTTP,
		)
	case expected.HTTP != nil:
		t.Errorf("Expected %s.HTTP %+v, got nil", path, expected.HTTP)
	case actual.HTTP != nil:
		t.Errorf("Expected %s.HTTP nil, got %+v", path, actual.HTTP)
	}

	switch {
	case expected.GRPC != nil && actual.GRPC != nil:
		assertEqualGRPCRequest(
			t,
			path+".GRPC",
			expected.GRPC,
			actual.GRPC,
		)
	case expected.GRPC != nil:
		t.Errorf("Expected %s.GRPC %+v, got nil", path, expected.GRPC)
	case actual.GRPC != nil:
		t.Errorf("Expected %s.GRPC nil, got %+v", path, actual.GRPC)
	}
}

func assertEqualTestExecutionResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionResponseSummary,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected %s nil, got %+v", path, actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected %s %+v, got nil", path, expected)
		return
	}

	switch {
	case expected.HTTP != nil && actual.HTTP != nil:
		assertEqualHTTPResponse(
			t,
			path+".HTTP",
			expected.HTTP,
			actual.HTTP,
		)
	case expected.HTTP != nil:
		t.Errorf("Expected %s.HTTP %+v, got nil", path, expected.HTTP)
	case actual.HTTP != nil:
		t.Errorf("Expected %s.HTTP nil, got %+v", path, actual.HTTP)
	}

	switch {
	case expected.GRPC != nil && actual.GRPC != nil:
		assertEqualGRPCResponse(
			t,
			path+".GRPC",
			expected.GRPC,
			actual.GRPC,
		)
	case expected.GRPC != nil:
		t.Errorf("Expected %s.GRPC %+v, got nil", path, expected.GRPC)
	case actual.GRPC != nil:
		t.Errorf("Expected %s.GRPC nil, got %+v", path, actual.GRPC)
	}
}

func assertEqualAllServiceObservedCalls(
	t *testing.T,
	path string,
	expected,
	actual []model.ServiceObservedCalls,
) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf(
			"Expected %s length %d, got %d",
			path,
			len(expected),
			len(actual),
		)
		return
	}

	for i := range expected {
		assertEqualServiceObservedCalls(
			t,
			path+"["+strconv.Itoa(i)+"]",
			&expected[i],
			&actual[i],
		)
	}
}

func assertEqualServiceObservedCalls(
	t *testing.T,
	path string,
	expected,
	actual *model.ServiceObservedCalls,
) {
	t.Helper()

	if expected.ServiceID != actual.ServiceID {
		t.Errorf("Expected %s.ServiceID %q, got %q", path, expected.ServiceID, actual.ServiceID)
	}

	if expected.Count != actual.Count {
		t.Errorf("Expected %s.Count %d, got %d", path, expected.Count, actual.Count)
	}

	assertEqualHTTPServiceObservedCall(
		t,
		path+".HTTP",
		expected.HTTP,
		actual.HTTP,
	)

	assertEqualGRPCServiceObservedCall(
		t,
		path+".GRPC",
		expected.GRPC,
		actual.GRPC,
	)
}

func assertEqualHTTPServiceObservedCall(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPServiceObservedCall,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected %s nil, got %+v", path, actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected %s %+v, got nil", path, expected)
		return
	}

	if expected.Method != actual.Method {
		t.Errorf("Expected %s.Method %q, got %q", path, expected.Method, actual.Method)
	}
	if expected.Path != actual.Path {
		t.Errorf("Expected %s.Path %q, got %q", path, expected.Path, actual.Path)
	}

	assertEqualIntPointer(
		t,
		path+".ResponseStatusCode",
		expected.ResponseStatusCode,
		actual.ResponseStatusCode,
	)
}

func assertEqualGRPCServiceObservedCall(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCServiceObservedCall,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected %s nil, got %+v", path, actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected %s %+v, got nil", path, expected)
		return
	}

	if expected.Service != actual.Service {
		t.Errorf("Expected %s.Service %q, got %q", path, expected.Service, actual.Service)
	}
	if expected.Method != actual.Method {
		t.Errorf("Expected %s.Method %q, got %q", path, expected.Method, actual.Method)
	}
	if expected.ResponseStatus != actual.ResponseStatus {
		t.Errorf("Expected %s.ResponseStatus %q, got %q", path, expected.ResponseStatus, actual.ResponseStatus)
	}
}

func assertEqualTestExecutionDeviations(
	t *testing.T,
	path string,
	expected,
	actual []model.TestExecutionDeviation,
) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %s length %d, got %d", path, len(expected), len(actual))
		return
	}

	for i := range expected {
		expectedDeviation := &expected[i]
		actualDeviation := &actual[i]
		itemPath := path + "[" + strconv.Itoa(i) + "]"

		if expectedDeviation.Field != actualDeviation.Field {
			t.Errorf(
				"Expected %s.Field %q, got %q",
				itemPath,
				expectedDeviation.Field,
				actualDeviation.Field,
			)
		}

		if expectedDeviation.Expected != actualDeviation.Expected {
			t.Errorf(
				"Expected %s.Expected %q, got %q",
				itemPath,
				expectedDeviation.Expected,
				actualDeviation.Expected,
			)
		}

		if expectedDeviation.Actual != actualDeviation.Actual {
			t.Errorf(
				"Expected %s.Actual %q, got %q",
				itemPath,
				expectedDeviation.Actual,
				actualDeviation.Actual,
			)
		}

		if expectedDeviation.Message != actualDeviation.Message {
			t.Errorf(
				"Expected %s.Message %q, got %q",
				itemPath,
				expectedDeviation.Message,
				actualDeviation.Message,
			)
		}
	}
}
