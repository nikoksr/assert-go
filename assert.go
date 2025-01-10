package assert

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//nolint:gochecknoglobals // read-only global
var activeConfig = Config{
	IncludeSource: true,
	ContextLines:  5, //nolint:gomnd,mnd // default value
}

// SetConfig sets the configuration for the assertion library.
func SetConfig(config Config) {
	activeConfig = config
}

// Assert panics if the condition is false. Gloabbly configurable via SetConfig.
//
// Assert is intended for critical checks that should always be active, regardless of the build configuration. Use
// assert.Debug for non-critical checks that should only be active during development.
//
// WARN: This assertion is live!
func Assert(condition bool, msg string, values ...any) {
	assert(condition, msg, values...)
}

// Assert panics if the condition is false. Configurable via SetConfig.
func assert(condition bool, msg string, values ...any) {
	if condition {
		return // Assertion met
	}

	// Skip 2 frames:
	// 1. this assert() function
	// 2. the Assert() function that called us
	_, file, line, _ := runtime.Caller(2) //nolint:mnd // Explained in comment

	// If values were provided for dumping
	numValues := len(values)
	if numValues%2 != 0 {
		values = append(values, "(MISSING)")
	}

	var dumpInfo string
	if numValues > 0 {
		dumpInfo = "\n\nRelevant values:\n"
		for i := 0; i < numValues; i += 2 {
			dumpInfo += fmt.Sprintf("  [%s]: %#v\n", values[i], values[i+1])
		}
	}

	// Get source context if enabled
	var sourceContext string
	if activeConfig.IncludeSource {
		sourceContext = getSourceContext(file, line, activeConfig.ContextLines)
	}

	err := AssertionError{
		Message:       msg + dumpInfo,
		File:          filepath.Base(file),
		Line:          line,
		SourceContext: sourceContext,
	}

	panic(err)
}

// getSourceContext reads the source file and returns lines around the failure.
func getSourceContext(file string, line int, contextLines int) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	start := max(1, line-contextLines)
	end := line + contextLines

	var lines []string
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= start && currentLine <= end {
			prefix := "  "
			if currentLine == line {
				prefix = "→ "
			}
			lines = append(lines, fmt.Sprintf("%s%4d | %s", prefix, currentLine, scanner.Text()))
		}

		if currentLine > end {
			break
		}

		currentLine++
	}

	return strings.Join(lines, "\n")
}
