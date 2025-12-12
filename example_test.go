package assert_test

import (
	"fmt"

	"github.com/nikoksr/assert-go"
)

// ExampleAssert demonstrates basic usage of the Assert function.
// This will panic if the condition is false.
func ExampleAssert() {
	// This assertion passes - no panic
	assert.Assert(true, "this should not panic")

	// In real code, you might check invariants like:
	// assert.Assert(user != nil, "user must be initialized before use",
	// 	"user_id", userID,
	// )

	fmt.Println("Assertion passed")
	// Output: Assertion passed
}

// ExampleAssert_withContext demonstrates using Assert with contextual values.
// The key-value pairs provide additional debugging information when assertions fail.
func ExampleAssert_withContext() {
	userID := "user_123"
	balance := 100.50
	status := "active"

	// This passes - demonstrating the API
	assert.Assert(status == "active", "user must be active",
		"user_id", userID,
		"balance", balance,
		"status", status,
	)

	fmt.Println("User validation passed")
	// Output: User validation passed
}

// ExampleSetConfig demonstrates configuring assertion behavior.
func ExampleSetConfig() {
	// Configure assertion behavior
	assert.SetConfig(assert.Config{
		// Enable source context in error messages
		IncludeSource: true,
		// Show 3 lines of context before and after the failing line
		ContextLines: 3,
	})

	// All subsequent assertions will use this configuration
	assert.Assert(true, "this uses the new configuration")

	fmt.Println("Configuration applied")
	// Output: Configuration applied
}

// ExampleDebug demonstrates the Debug assertion that can be enabled with build tags.
// Debug assertions are disabled by default and only active when built with -tags assertdebug.
func ExampleDebug() {
	// This assertion is only evaluated when built with: go test -tags assertdebug
	assert.Debug(true, "this check only runs in debug builds")

	// Use Debug for non-critical checks during development:
	// assert.Debug(len(cache) < 10000, "cache getting large",
	// 	"size", len(cache),
	// )

	fmt.Println("Debug assertion evaluated")
	// Output: Debug assertion evaluated
}
