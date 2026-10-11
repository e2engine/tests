package util

import (
	"testing"

	"github.com/e2engine/core/model"
	"github.com/e2engine/core/model/matcher"
)

func assertEqualHTTPRequest(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPRequest,
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

	if expected.BodyJSON != actual.BodyJSON {
		t.Errorf(
			"Expected %s.BodyJSON %q, got %q",
			path,
			expected.BodyJSON,
			actual.BodyJSON,
		)
	}
}

func assertEqualHTTPRequestPattern(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPRequestPattern,
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

	assertEqualStringMatcherSliceMap(
		t,
		path+".Headers",
		expected.Headers,
		actual.Headers,
	)

	assertEqualPtrStringMatcher(
		t,
		path+".Body",
		expected.Body,
		actual.Body,
	)

	assertEqualPtrStringMatcher(
		t,
		path+".BodyJSON",
		expected.BodyJSON,
		actual.BodyJSON,
	)
}

func assertEqualPtrStringMatcher(
	t *testing.T,
	path string,
	expected,
	actual *matcher.String,
) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected == nil {
		t.Errorf("Expected %s nil, got %s", path, *actual)
	} else if actual == nil {
		t.Errorf("Expected %s %s, got nil", path, *expected)
	} else if *expected != *actual {
		t.Errorf(
			"Expected %s %s, got %s",
			path,
			*expected,
			*actual,
		)
	}
}

func assertEqualHTTPResponsePattern(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPResponsePattern,
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

	if expected.StatusCode != actual.StatusCode {
		t.Errorf(
			"Expected %s.StatusCode %d, got %d",
			path,
			expected.StatusCode,
			actual.StatusCode,
		)
	}

	assertEqualStringMatcherSliceMap(
		t,
		path+".Headers",
		expected.Headers,
		actual.Headers,
	)

	assertEqualPtrStringMatcher(
		t,
		path+".Body",
		expected.Body,
		actual.Body,
	)

	assertEqualPtrStringMatcher(
		t,
		path+".BodyJSON",
		expected.BodyJSON,
		actual.BodyJSON,
	)
}

func assertEqualHTTPResponse(
	t *testing.T,
	path string,
	expected,
	actual *model.HTTPResponse,
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

	if expected.StatusCode != actual.StatusCode {
		t.Errorf(
			"Expected %s.StatusCode %d, got %d",
			path,
			expected.StatusCode,
			actual.StatusCode,
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

	if expected.BodyJSON != actual.BodyJSON {
		t.Errorf(
			"Expected %s.BodyJSON %q, got %q",
			path,
			expected.BodyJSON,
			actual.BodyJSON,
		)
	}
}
