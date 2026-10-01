package util

import (
	"testing"

	"github.com/e2engine/core/model"
)

func AssertEqualTestSuites(t *testing.T, expected, actual *model.TestSuite) {
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

	assertEqualTestSuiteSpec(t, &expected.Spec, &actual.Spec)
}

func assertEqualTestSuiteSpec(
	t *testing.T,
	expected,
	actual *model.TestSuiteSpec,
) {
	t.Helper()

	assertEqualStringSlice(t, "Selectors.IDs", expected.Selectors.IDs, actual.Selectors.IDs)
	assertEqualStringSlice(t, "Selectors.Names", expected.Selectors.Names, actual.Selectors.Names)
	assertEqualStringSlice(t, "Selectors.Tags", expected.Selectors.Tags, actual.Selectors.Tags)
}
