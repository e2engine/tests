package util

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/e2engine/core/model"
)

func AssertEqualEnvironments(t *testing.T, expected, actual *model.Environment) {
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

	assertEqualEnvironmentSpec(t, &expected.Spec, &actual.Spec)
}

func assertEqualEnvironmentSpec(
	t *testing.T,
	expected,
	actual *model.EnvironmentSpec,
) {
	t.Helper()

	if len(expected.Services) != len(actual.Services) {
		t.Errorf(
			"Expected %d services, got %d",
			len(expected.Services),
			len(actual.Services),
		)
		return
	}

	for i := range expected.Services {
		assertEqualServiceSpec(
			t,
			i,
			&expected.Services[i],
			&actual.Services[i],
		)
	}
}

func assertEqualServiceSpec(
	t *testing.T,
	index int,
	expected,
	actual *model.ServiceSpec,
) {
	t.Helper()

	prefix := func(field string) string {
		return "Spec.Services[" + strconv.Itoa(index) + "]." + field
	}

	if expected.ID != actual.ID {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("ID"),
			expected.ID,
			actual.ID,
		)
	}

	if expected.Kind != actual.Kind {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("Kind"),
			expected.Kind,
			actual.Kind,
		)
	}

	if expected.Mode != actual.Mode {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("Mode"),
			expected.Mode,
			actual.Mode,
		)
	}

	if expected.Address != actual.Address {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("Address"),
			expected.Address,
			actual.Address,
		)
	}

	if expected.HTTPTarget != actual.HTTPTarget {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("HTTPTarget"),
			expected.HTTPTarget,
			actual.HTTPTarget,
		)
	}

	if expected.GRPCTarget != actual.GRPCTarget {
		t.Errorf(
			"Expected %s %q, got %q",
			prefix("GRPCTarget"),
			expected.GRPCTarget,
			actual.GRPCTarget,
		)
	}

	assertEqualServiceDefaults(
		t,
		prefix("RequestDefaults"),
		expected.RequestDefaults,
		actual.RequestDefaults,
	)

	assertEqualGRPCProtoSpec(
		t,
		prefix("Proto"),
		expected.Proto,
		actual.Proto,
	)

	if len(expected.Fixtures) != len(actual.Fixtures) {
		t.Errorf(
			"Expected %s length %d, got %d",
			prefix("Fixtures"),
			len(expected.Fixtures),
			len(actual.Fixtures),
		)
		return
	}

	for i := range expected.Fixtures {
		assertEqualFixtureSpec(
			t,
			prefix("Fixtures")+"["+strconv.Itoa(i)+"]",
			&expected.Fixtures[i],
			&actual.Fixtures[i],
		)
	}
}

func assertEqualServiceDefaults(
	t *testing.T,
	path string,
	expected,
	actual *model.ServiceDefaults,
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

	assertEqualStringSliceMap(
		t,
		path+".Headers",
		expected.Headers,
		actual.Headers,
	)
}

func assertEqualGRPCProtoSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCProtoSpec,
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

	assertEqualGRPCExternalProtoSpec(
		t,
		path+".External",
		expected.External,
		actual.External,
	)

	assertEqualGRPCInternalProtoSpec(
		t,
		path+".Internal",
		expected.Internal,
		actual.Internal,
	)
}

func assertEqualGRPCExternalProtoSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCExternalProtoSpec,
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

	if expected.BufModule != actual.BufModule {
		t.Errorf(
			"Expected %s.BufModule %q, got %q",
			path,
			expected.BufModule,
			actual.BufModule,
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
}

func assertEqualGRPCInternalProtoSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCInternalProtoSpec,
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

	if expected.Package != actual.Package {
		t.Errorf(
			"Expected %s.Package %q, got %q",
			path,
			expected.Package,
			actual.Package,
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

	if len(expected.Methods) != len(actual.Methods) {
		t.Errorf(
			"Expected %s.Methods length %d, got %d",
			path,
			len(expected.Methods),
			len(actual.Methods),
		)
		return
	}

	for i := range expected.Methods {
		assertEqualGRPCMethodSpec(
			t,
			path+".Methods["+strconv.Itoa(i)+"]",
			&expected.Methods[i],
			&actual.Methods[i],
		)
	}
}

func assertEqualGRPCMethodSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCMethodSpec,
) {
	t.Helper()

	if expected.Name != actual.Name {
		t.Errorf(
			"Expected %s.Name %q, got %q",
			path,
			expected.Name,
			actual.Name,
		)
	}

	assertEqualGRPCMessageSpec(
		t,
		path+".Request",
		&expected.Request,
		&actual.Request,
	)

	assertEqualGRPCMessageSpec(
		t,
		path+".Response",
		&expected.Response,
		&actual.Response,
	)
}

func assertEqualGRPCMessageSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCMessageSpec,
) {
	t.Helper()

	if len(expected.Fields) != len(actual.Fields) {
		t.Errorf(
			"Expected %s.Fields length %d, got %d",
			path,
			len(expected.Fields),
			len(actual.Fields),
		)
		return
	}

	for i := range expected.Fields {
		assertEqualGRPCFieldSpec(
			t,
			path+".Fields["+strconv.Itoa(i)+"]",
			&expected.Fields[i],
			&actual.Fields[i],
		)
	}
}

func assertEqualGRPCFieldSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCFieldSpec,
) {
	t.Helper()

	if expected.Name != actual.Name {
		t.Errorf(
			"Expected %s.Name %q, got %q",
			path,
			expected.Name,
			actual.Name,
		)
	}

	if expected.Type != actual.Type {
		t.Errorf(
			"Expected %s.Type %q, got %q",
			path,
			expected.Type,
			actual.Type,
		)
	}

	if expected.Repeated != actual.Repeated {
		t.Errorf(
			"Expected %s.Repeated %t, got %t",
			path,
			expected.Repeated,
			actual.Repeated,
		)
	}
}

func assertEqualFixtureSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.FixtureSpec,
) {
	t.Helper()

	assertEqualFixtureWhen(
		t,
		path+".When",
		&expected.When,
		&actual.When,
	)

	assertEqualFixtureThen(
		t,
		path+".Then",
		&expected.Then,
		&actual.Then,
	)
}

func assertEqualFixtureWhen(
	t *testing.T,
	path string,
	expected,
	actual *model.FixtureWhen,
) {
	t.Helper()

	assertEqualHTTPFixtureWhen(
		t,
		path+".HTTP",
		expected.HTTP,
		actual.HTTP,
	)

	assertEqualGRPCFixtureWhen(
		t,
		path+".GRPC",
		expected.GRPC,
		actual.GRPC,
	)
}

func assertEqualFixtureThen(
	t *testing.T,
	path string,
	expected,
	actual *model.FixtureThen,
) {
	t.Helper()

	assertEqualHTTPFixtureThen(
		t,
		path+".HTTP",
		expected.HTTP,
		actual.HTTP,
	)

	assertEqualGRPCFixtureThen(
		t,
		path+".GRPC",
		expected.GRPC,
		actual.GRPC,
	)
}

func assertEqualHTTPFixtureWhen(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPFixtureWhen,
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

func assertEqualHTTPFixtureThen(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPFixtureThen,
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

func assertEqualGRPCFixtureWhen(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCFixtureWhen,
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

	assertEqualJSONValue(
		t,
		path+".Message",
		expected.Message,
		actual.Message,
	)
}

func assertEqualGRPCFixtureThen(
	t *testing.T,
	path string,
	expected,
	actual *model.GRPCFixtureThen,
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

func assertEqualStringSliceMap(
	t *testing.T,
	path string,
	expected,
	actual map[string][]string,
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

	for key, expectedValues := range expected {
		actualValues, ok := actual[key]
		if !ok {
			t.Errorf("Expected %s[%q] to exist", path, key)
			continue
		}

		if len(expectedValues) != len(actualValues) {
			t.Errorf(
				"Expected %s[%q] length %d, got %d",
				path,
				key,
				len(expectedValues),
				len(actualValues),
			)
			continue
		}

		if key == "Date" {
			continue // TODO: replace with check.
		}
		for i := range expectedValues {
			if expectedValues[i] != actualValues[i] {
				t.Errorf(
					"Expected %s[%q][%d] %q, got %q",
					path,
					key,
					i,
					expectedValues[i],
					actualValues[i],
				)
			}
		}
	}
}

func assertEqualJSONValue(
	t *testing.T,
	path string,
	expected,
	actual any,
) {
	t.Helper()

	expectedJSON, expectedErr := json.Marshal(expected)
	actualJSON, actualErr := json.Marshal(actual)

	if expectedErr != nil {
		t.Fatalf(
			"Failed to marshal expected %s: %v",
			path,
			expectedErr,
		)
	}

	if actualErr != nil {
		t.Fatalf(
			"Failed to marshal actual %s: %v",
			path,
			actualErr,
		)
	}

	if !bytes.Equal(expectedJSON, actualJSON) {
		t.Errorf(
			"Expected %s %s, got %s",
			path,
			expectedJSON,
			actualJSON,
		)
	}
}
