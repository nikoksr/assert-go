//go:build !assertdebug

package assert

// Debug is a no-op. The `assertdebug` build tag is not set.
//
// To learn more about build tags, see https://pkg.go.dev/go/build#hdr-Build_Constraints.
func Debug(_ bool, _ string, _ ...any) {}
