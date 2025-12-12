package assert

import "testing"

// BenchmarkAssert_Success benchmarks the successful assertion path.
// This is the hot path - assertions that pass should be extremely cheap.
func BenchmarkAssert_Success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Assert(true, "never fails")
	}
}

// BenchmarkAssert_SuccessWithValues benchmarks successful assertions with contextual values.
// Even though values are provided, they should not be evaluated when the assertion passes.
func BenchmarkAssert_SuccessWithValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Assert(true, "never fails",
			"key1", "value1",
			"key2", 42,
			"key3", struct{ Field string }{Field: "test"},
		)
	}
}

// BenchmarkAssert_WithSourceContext benchmarks assertions with source context enabled.
func BenchmarkAssert_WithSourceContext(b *testing.B) {
	// Save and restore original config
	originalConfig := activeConfig
	defer func() {
		SetConfig(originalConfig)
	}()

	SetConfig(Config{IncludeSource: true, ContextLines: 5})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Assert(true, "with context")
	}
}

// BenchmarkAssert_WithoutSourceContext benchmarks assertions with source context disabled.
func BenchmarkAssert_WithoutSourceContext(b *testing.B) {
	// Save and restore original config
	originalConfig := activeConfig
	defer func() {
		SetConfig(originalConfig)
	}()

	SetConfig(Config{IncludeSource: false})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Assert(true, "no context")
	}
}

// BenchmarkDebug_Success benchmarks the Debug assertion (disabled by default).
func BenchmarkDebug_Success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debug(true, "debug assertion")
	}
}

// BenchmarkDebug_SuccessWithValues benchmarks Debug assertions with values.
func BenchmarkDebug_SuccessWithValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debug(true, "debug assertion",
			"key1", "value1",
			"key2", 42,
		)
	}
}
