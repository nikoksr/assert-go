//go:build noassert
// +build noassert

package assert

// SetConfig is a no-op when assertions are disabled
func SetConfig(config Config) {}

// Assert is a no-op when assertions are disabled
func Assert(condition bool, msg string, values ...any) {}
