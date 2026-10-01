package util

import (
	"strconv"
	"strings"
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
		&expected.Request,
		&actual.Request,
	)

	assertEqualTestExecutionResponseSummary(
		t,
		"Summary.Response",
		&expected.Response,
		&actual.Response,
	)

	assertEqualTestExecutionExpectSummary(
		t,
		"Summary.Expect",
		&expected.Expect,
		&actual.Expect,
	)

	assertEqualTestExecutionCalls(
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

	switch {
	case expected.HTTP != nil && actual.HTTP != nil:
		assertEqualTestExecutionHTTPRequestSummary(
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
		assertEqualTestExecutionGRPCRequestSummary(
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

func assertEqualTestExecutionHTTPRequestSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPRequestSummary,
) {
	t.Helper()

	if expected.Method != actual.Method {
		t.Errorf("Expected %s.Method %q, got %q", path, expected.Method, actual.Method)
	}
	if expected.URL != actual.URL {
		t.Errorf("Expected %s.URL %q, got %q", path, expected.URL, actual.URL)
	}

	assertEqualStringSliceMap(t, path+".Headers", expected.Headers, actual.Headers)

	if expected.BodyJSON != actual.BodyJSON {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
}

func assertEqualTestExecutionGRPCRequestSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCRequestSummary,
) {
	t.Helper()

	if expected.Service != actual.Service {
		t.Errorf("Expected %s.Service %q, got %q", path, expected.Service, actual.Service)
	}
	if expected.Method != actual.Method {
		t.Errorf("Expected %s.Method %q, got %q", path, expected.Method, actual.Method)
	}

	assertEqualStringSliceMap(t, path+".Metadata", expected.Metadata, actual.Metadata)
	assertEqualJSONValue(t, path+".Message", expected.Message, actual.Message)
}

func assertEqualTestExecutionResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionResponseSummary,
) {
	t.Helper()

	switch {
	case expected.HTTP != nil && actual.HTTP != nil:
		assertEqualTestExecutionHTTPResponseSummary(
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
		assertEqualTestExecutionGRPCResponseSummary(
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

func assertEqualTestExecutionHTTPResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPResponseSummary,
) {
	t.Helper()

	if expected.StatusCode != actual.StatusCode {
		t.Errorf("Expected %s.StatusCode %d, got %d", path, expected.StatusCode, actual.StatusCode)
	}
	if expected.BodyJSON != actual.BodyJSON {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
	if expected.BodyText != strings.TrimSpace(actual.BodyText) {
		t.Errorf("Expected %s.BodyText %q, got %q", path, expected.BodyText, actual.BodyText)
	}
}

func assertEqualTestExecutionGRPCResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCResponseSummary,
) {
	t.Helper()

	if expected.Status != actual.Status {
		t.Errorf("Expected %s.Status %q, got %q", path, expected.Status, actual.Status)
	}

	assertEqualStringSliceMap(t, path+".Metadata", expected.Metadata, actual.Metadata)
	assertEqualJSONValue(t, path+".Message", expected.Message, actual.Message)
}

func assertEqualTestExecutionExpectSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionExpectSummary,
) {
	t.Helper()

	switch {
	case expected.HTTP != nil && actual.HTTP != nil:
		assertEqualTestExecutionHTTPExpectSummary(
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
		assertEqualTestExecutionGRPCExpectSummary(
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

	if len(expected.Calls) != len(actual.Calls) {
		t.Errorf(
			"Expected %s.Calls length %d, got %d",
			path,
			len(expected.Calls),
			len(actual.Calls),
		)
		return
	}

	for i := range expected.Calls {
		assertEqualTestExecutionCallExpectationSummary(
			t,
			path+".Calls["+strconv.Itoa(i)+"]",
			&expected.Calls[i],
			&actual.Calls[i],
		)
	}
}

func assertEqualTestExecutionHTTPExpectSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPExpectSummary,
) {
	t.Helper()

	if expected.StatusCode != actual.StatusCode {
		t.Errorf("Expected %s.StatusCode %d, got %d", path, expected.StatusCode, actual.StatusCode)
	}
	if expected.BodyJSON != actual.BodyJSON {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
}

func assertEqualTestExecutionGRPCExpectSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCExpectSummary,
) {
	t.Helper()

	if expected.Status != actual.Status {
		t.Errorf("Expected %s.Status %q, got %q", path, expected.Status, actual.Status)
	}
	assertEqualJSONValue(t, path+".Message", expected.Message, actual.Message)
}

func assertEqualTestExecutionCallExpectationSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionCallExpectationSummary,
) {
	t.Helper()

	if expected.ServiceID != actual.ServiceID {
		t.Errorf("Expected %s.ServiceID %q, got %q", path, expected.ServiceID, actual.ServiceID)
	}

	assertEqualIntPointer(t, path+".Count", expected.Count, actual.Count)
	assertEqualTestExecutionHTTPCallExpectationSummary(
		t,
		path+".HTTP",
		expected.HTTP,
		actual.HTTP,
	)
}

func assertEqualTestExecutionHTTPCallExpectationSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPCallExpectationSummary,
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

	assertEqualStringSliceMap(t, path+".Query", expected.Query, actual.Query)
	assertEqualStringSliceMap(t, path+".Headers", expected.Headers, actual.Headers)

	if expected.BodyJSON != strings.TrimSpace(actual.BodyJSON) {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
}

func assertEqualTestExecutionCalls(
	t *testing.T,
	path string,
	expected,
	actual []model.TestExecutionCallSummary,
) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %s length %d, got %d", path, len(expected), len(actual))
		return
	}

	for i := range expected {
		assertEqualTestExecutionCallSummary(
			t,
			path+"["+strconv.Itoa(i)+"]",
			&expected[i],
			&actual[i],
		)
	}
}

func assertEqualTestExecutionCallSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionCallSummary,
) {
	t.Helper()

	if expected.ServiceID != actual.ServiceID {
		t.Errorf("Expected %s.ServiceID %q, got %q", path, expected.ServiceID, actual.ServiceID)
	}

	assertEqualTestExecutionHTTPCallSummary(
		t,
		path+".HTTP",
		expected.HTTP,
		actual.HTTP,
	)

	assertEqualTestExecutionGRPCCallSummary(
		t,
		path+".GRPC",
		expected.GRPC,
		actual.GRPC,
	)
}

func assertEqualTestExecutionGRPCCallSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCCallSummary,
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

	assertEqualTestExecutionGRPCCallRequestSummary(
		t,
		path+".Request",
		&expected.Request,
		&actual.Request,
	)

	assertEqualTestExecutionGRPCCallResponseSummary(
		t,
		path+".Response",
		&expected.Response,
		&actual.Response,
	)
}

func assertEqualTestExecutionGRPCCallRequestSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCCallRequestSummary,
) {
	t.Helper()

	if expected.RPC != actual.RPC {
		t.Errorf("Expected %s.RPC %q, got %q", path, expected.RPC, actual.RPC)
	}

	assertEqualStringSliceMap(
		t,
		path+".Metadata",
		expected.Metadata,
		actual.Metadata,
	)

	assertEqualJSONValue(
		t,
		path+".Message",
		expected.Message,
		actual.Message,
	)
}

func assertEqualTestExecutionGRPCCallResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionGRPCCallResponseSummary,
) {
	t.Helper()

	if expected.Status != actual.Status {
		t.Errorf("Expected %s.Status %q, got %q", path, expected.Status, actual.Status)
	}

	assertEqualStringSliceMap(
		t,
		path+".Metadata",
		expected.Metadata,
		actual.Metadata,
	)

	assertEqualJSONValue(
		t,
		path+".Message",
		expected.Message,
		actual.Message,
	)
}

func assertEqualTestExecutionHTTPCallSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPCallSummary,
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

	assertEqualTestExecutionHTTPCallRequestSummary(
		t,
		path+".Request",
		&expected.Request,
		&actual.Request,
	)

	assertEqualTestExecutionHTTPCallResponseSummary(
		t,
		path+".Response",
		&expected.Response,
		&actual.Response,
	)
}

func assertEqualTestExecutionHTTPCallRequestSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPCallRequestSummary,
) {
	t.Helper()

	if expected.Method != actual.Method {
		t.Errorf("Expected %s.Method %q, got %q", path, expected.Method, actual.Method)
	}
	if expected.Path != actual.Path {
		t.Errorf("Expected %s.Path %q, got %q", path, expected.Path, actual.Path)
	}

	assertEqualStringSliceMap(t, path+".Query", expected.Query, actual.Query)
	assertEqualStringSliceMap(t, path+".Headers", expected.Headers, actual.Headers)

	if expected.BodyJSON != strings.TrimSpace(actual.BodyJSON) {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
}

func assertEqualTestExecutionHTTPCallResponseSummary(
	t *testing.T,
	path string,
	expected,
	actual *model.TestExecutionHTTPCallResponseSummary,
) {
	t.Helper()

	if expected.StatusCode != actual.StatusCode {
		t.Errorf("Expected %s.StatusCode %d, got %d", path, expected.StatusCode, actual.StatusCode)
	}

	assertEqualStringSliceMap(t, path+".Headers", expected.Headers, actual.Headers)

	if expected.BodyJSON != strings.TrimSpace(actual.BodyJSON) {
		t.Errorf("Expected %s.BodyJSON %q, got %q", path, expected.BodyJSON, actual.BodyJSON)
	}
	if expected.BodyText != strings.TrimSpace(actual.BodyText) {
		t.Errorf("Expected %s.BodyText %q, got %q", path, expected.BodyText, actual.BodyText)
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
