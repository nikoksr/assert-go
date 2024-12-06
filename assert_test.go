package assert

import (
	"strings"
	"testing"
)

func TestAssert_SuccessCase(t *testing.T) {
	t.Parallel()
	Assert(true, "this should not panic")
}

func TestAssert_BasicFailure(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but got none")
		}

		err, ok := r.(AssertionError)
		if !ok {
			t.Fatalf("expected AssertionError but got %T", r)
		}

		// Check basic error properties
		if err.Message != "basic failure message" {
			t.Errorf("expected message 'basic failure message' but got '%s'", err.Message)
		}

		// Only verify that line number exists and is positive
		if err.Line <= 0 {
			t.Error("expected positive line number")
		}

		// Only verify file name format, not specific line
		if !strings.HasSuffix(err.File, "_test.go") {
			t.Errorf("expected file name to end with _test.go, got %s", err.File)
		}

		// Verify source context exists and contains the failure line
		if !strings.Contains(err.SourceContext, "Assert(false, \"basic failure message\")") {
			t.Error("expected source context to contain the failing assertion")
		}
	}()

	Assert(false, "basic failure message")
}

func TestAssert_WithValues(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field string
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but got none")
		}

		err, ok := r.(AssertionError)
		if !ok {
			t.Fatalf("expected AssertionError but got %T", r)
		}

		expectedValues := []string{
			"[string_key]: \"string_value\"",
			"[int_key]: 42",
			"[struct_key]:",
		}

		for _, expected := range expectedValues {
			if !strings.Contains(err.Message, expected) {
				t.Errorf("expected message to contain '%s'", expected)
			}
		}

		// Verify source context exists and contains the failure line
		if !strings.Contains(err.SourceContext, "Assert(false, \"failure with values\"") {
			t.Error("expected source context to contain the failing assertion")
		}
	}()

	Assert(false, "failure with values",
		"string_key", "string_value",
		"int_key", 42,
		"struct_key", testStruct{Field: "value"},
	)
}

func TestAssert_OddNumberOfValues(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but got none")
		}

		err, ok := r.(AssertionError)
		if !ok {
			t.Fatalf("expected AssertionError but got %T", r)
		}

		if !strings.Contains(err.Message, "(MISSING)") {
			t.Error("expected message to contain (MISSING) for odd number of values")
		}

		// Verify source context exists
		if err.SourceContext == "" {
			t.Error("expected non-empty source context")
		}
	}()

	Assert(false, "odd values",
		"key1", "value1",
		"key2", // Missing value
	)
}

func TestAssertionError_Error(t *testing.T) {
	t.Parallel()

	err := AssertionError{
		Message: "test message",
		File:    "test_file.go",
		Line:    42,
		SourceContext: "   41 | func TestExample(t *testing.T) {\n" +
			"→  42 | \tAssert(false, \"test message\")\n" +
			"   43 | }",
	}

	errStr := err.Error()

	expectedParts := []string{
		"test_file.go:42",
		"test message",
		"Source context:",
		"→  42",
	}

	for _, part := range expectedParts {
		if !strings.Contains(errStr, part) {
			t.Errorf("expected error string to contain '%s'", part)
		}
	}
}

func TestAssert_NilValues(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but got none")
		}

		err, ok := r.(AssertionError)
		if !ok {
			t.Fatalf("expected AssertionError but got %T", r)
		}

		if !strings.Contains(err.Message, "[nil_key]: <nil>") {
			t.Error("expected message to handle nil values correctly")
		}

		// Verify source context exists
		if err.SourceContext == "" {
			t.Error("expected non-empty source context")
		}
	}()

	Assert(false, "nil value test",
		"nil_key", nil,
	)
}

func TestAssert_EmptyValues(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic but got none")
		}

		err, ok := r.(AssertionError)
		if !ok {
			t.Fatalf("expected AssertionError but got %T", r)
		}

		expectedValues := []string{
			"[empty_string]: \"\"",
			"[empty_slice]: []string{}",
			"[empty_map]: map[string]int{}",
		}

		for _, expected := range expectedValues {
			if !strings.Contains(err.Message, expected) {
				t.Errorf("expected message to contain '%s'", expected)
			}
		}

		// Verify source context exists
		if err.SourceContext == "" {
			t.Error("expected non-empty source context")
		}
	}()

	Assert(false, "empty values test",
		"empty_string", "",
		"empty_slice", []string{},
		"empty_map", map[string]int{},
	)
}