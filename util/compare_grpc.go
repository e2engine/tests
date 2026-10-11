package util

import (
	"testing"

	"github.com/e2engine/core/model"
)

func assertEqualGRPCRequest(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCRequest,
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

func assertEqualGRPCRequestPattern(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCRequestPattern,
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

	assertEqualStringMatcherSliceMap(
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

func assertEqualGRPCResponsePattern(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCResponsePattern,
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

	assertEqualStringMatcherSliceMap(
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

func assertEqualGRPCResponse(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCResponse,
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
