package util

import (
	"strconv"
	"testing"

	"github.com/e2engine/core/model"
)

func AssertEqualTests(t *testing.T, expected, actual *model.Test) {
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
	if expected.Name != actual.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, actual.Name)
	}
	if expected.Description != actual.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, actual.Description)
	}
	if !expected.CreatedAt.Equal(actual.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", expected.CreatedAt, actual.CreatedAt)
	}
	if !expected.UpdatedAt.Equal(actual.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", expected.UpdatedAt, actual.UpdatedAt)
	}

	assertEqualTestSpec(t, &expected.Spec, &actual.Spec)
}

func assertEqualTestSpec(
	t *testing.T,
	expected,
	actual *model.TestSpec,
) {
	t.Helper()

	assertEqualStringSlice(t, "Spec.Tags", expected.Tags, actual.Tags)
	assertEqualRequestSpec(t, "Spec.Request", &expected.Request, &actual.Request)
	assertEqualExpectSpec(t, "Spec.Expect", &expected.Expect, &actual.Expect)
}

func assertEqualRequestSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.RequestSpec,
) {
	t.Helper()

	assertEqualHTTPRequestSpec(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCRequestSpec(t, path+".GRPC", expected.GRPC, actual.GRPC)
}

func assertEqualExpectSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.ExpectSpec,
) {
	t.Helper()

	assertEqualHTTPExpectSpec(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCExpectSpec(t, path+".GRPC", expected.GRPC, actual.GRPC)

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
		assertEqualCallExpectation(
			t,
			path+".Calls["+strconv.Itoa(i)+"]",
			&expected.Calls[i],
			&actual.Calls[i],
		)
	}
}

func assertEqualCallExpectation(
	t *testing.T,
	path string,
	expected,
	actual *model.CallExpectation,
) {
	t.Helper()

	if expected.ServiceID != actual.ServiceID {
		t.Errorf(
			"Expected %s.ServiceID %q, got %q",
			path,
			expected.ServiceID,
			actual.ServiceID,
		)
	}

	assertEqualIntPointer(t, path+".Count", expected.Count, actual.Count)
	assertEqualHTTPCallExpectation(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCCallExpectation(t, path+".GRPC", expected.GRPC, actual.GRPC)
}

func assertEqualHTTPCallExpectation(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPCallExpectation,
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
		t.Errorf(
			"Expected %s.Method %q, got %q",
			path,
			expected.Method,
			actual.Method,
		)
	}

	if expected.Path != actual.Path {
		t.Errorf(
			"Expected %s.Path %q, got %q",
			path,
			expected.Path,
			actual.Path,
		)
	}

	assertEqualStringSliceMap(
		t,
		path+".Query",
		expected.Query,
		actual.Query,
	)

	assertEqualStringSliceMap(
		t,
		path+".Headers",
		expected.Headers,
		actual.Headers,
	)

	if expected.Body != actual.Body {
		t.Errorf(
			"Expected %s.Body %q, got %q",
			path,
			expected.Body,
			actual.Body,
		)
	}
}

func assertEqualGRPCCallExpectation(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCCallExpectation,
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
		t.Errorf(
			"Expected %s.Service %q, got %q",
			path,
			expected.Service,
			actual.Service,
		)
	}

	if expected.Method != actual.Method {
		t.Errorf(
			"Expected %s.Method %q, got %q",
			path,
			expected.Method,
			actual.Method,
		)
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

func assertEqualHTTPRequestSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPRequestSpec,
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
		t.Errorf(
			"Expected %s.Method %q, got %q",
			path,
			expected.Method,
			actual.Method,
		)
	}

	if expected.URL != actual.URL {
		t.Errorf(
			"Expected %s.URL %q, got %q",
			path,
			expected.URL,
			actual.URL,
		)
	}

	assertEqualStringMap(
		t,
		path+".Query",
		expected.Query,
		actual.Query,
	)

	assertEqualStringSliceMap(
		t,
		path+".Headers",
		expected.Headers,
		actual.Headers,
	)

	if expected.Body != actual.Body {
		t.Errorf(
			"Expected %s.Body %q, got %q",
			path,
			expected.Body,
			actual.Body,
		)
	}
}

func assertEqualHTTPExpectSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPExpectSpec,
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

	if expected.Status != actual.Status {
		t.Errorf(
			"Expected %s.Status %d, got %d",
			path,
			expected.Status,
			actual.Status,
		)
	}

	if expected.Body != actual.Body {
		t.Errorf(
			"Expected %s.Body %q, got %q",
			path,
			expected.Body,
			actual.Body,
		)
	}
}

func assertEqualGRPCRequestSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCRequestSpec,
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

	if expected.Target != actual.Target {
		t.Errorf(
			"Expected %s.Target %q, got %q",
			path,
			expected.Target,
			actual.Target,
		)
	}

	if expected.Service != actual.Service {
		t.Errorf(
			"Expected %s.Service %q, got %q",
			path,
			expected.Service,
			actual.Service,
		)
	}

	if expected.Method != actual.Method {
		t.Errorf(
			"Expected %s.Method %q, got %q",
			path,
			expected.Method,
			actual.Method,
		)
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

func assertEqualGRPCExpectSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCExpectSpec,
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

	if expected.Status != actual.Status {
		t.Errorf(
			"Expected %s.Status %q, got %q",
			path,
			expected.Status,
			actual.Status,
		)
	}

	assertEqualJSONValue(
		t,
		path+".Message",
		expected.Message,
		actual.Message,
	)
}

func assertEqualIntPointer(
	t *testing.T,
	path string,
	expected,
	actual *int,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("Expected %s nil, got %d", path, *actual)
		return
	}
	if actual == nil {
		t.Errorf("Expected %s %d, got nil", path, *expected)
		return
	}

	if *expected != *actual {
		t.Errorf(
			"Expected %s %d, got %d",
			path,
			*expected,
			*actual,
		)
	}
}

func assertEqualStringSlice(
	t *testing.T,
	path string,
	expected,
	actual []string,
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
		if expected[i] != actual[i] {
			t.Errorf(
				"Expected %s[%d] %q, got %q",
				path,
				i,
				expected[i],
				actual[i],
			)
		}
	}
}

func assertEqualStringMap(
	t *testing.T,
	path string,
	expected,
	actual map[string]string,
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

	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok {
			t.Errorf(
				"Expected %s[%q] to exist",
				path,
				key,
			)
			continue
		}

		if expectedValue != actualValue {
			t.Errorf(
				"Expected %s[%q] %q, got %q",
				path,
				key,
				expectedValue,
				actualValue,
			)
		}
	}
}
