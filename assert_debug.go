//go:build assertdebug
// +build assertdebug

package assert

// Debug panics if the condition is false. Globally configurable via SetConfig.
//
// Debug is intended for non-critical checks that should only be active during development. For critical checks that
// should always be active, regardless of the build configuration, use assert.Assert instead.
//
// WARN: Under the current build configuration, this assertion is enabled.
func Debug(condition bool, msg string, values ...any) {
	assert(condition, msg, values...)
}
