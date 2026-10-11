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

	assertEqualHTTPRequest(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCRequest(t, path+".GRPC", expected.GRPC, actual.GRPC)
}

func assertEqualExpectSpec(
	t *testing.T,
	path string,
	expected,
	actual *model.ExpectSpec,
) {
	t.Helper()

	assertEqualHTTPResponsePattern(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCResponsePattern(t, path+".GRPC", expected.GRPC, actual.GRPC)

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
		assertEqualServiceExpectedCalls(
			t,
			path+".Calls["+strconv.Itoa(i)+"]",
			&expected.Calls[i],
			&actual.Calls[i],
		)
	}
}

func assertEqualServiceExpectedCalls(
	t *testing.T,
	path string,
	expected,
	actual *model.ServiceExpectedCalls,
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
	assertEqualHTTPRequestPattern(t, path+".HTTP", expected.HTTP, actual.HTTP)
	assertEqualGRPCRequestPattern(t, path+".GRPC", expected.GRPC, actual.GRPC)
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
