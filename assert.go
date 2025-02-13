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
	// We tell assert() to skip 2 frames here:
	//  1. The assert() function itself
	//  2. This Assert() function that calls assert()
	assert(condition, msg, 2, values...) //nolint:mnd // Explained in comment
}

// Assert panics if the condition is false. Configurable via SetConfig.
// skipFrames is the number of stack frames to skip when getting the source context. 
func assert(condition bool, msg string, skipFrames int, values ...any) {
	if condition {
		return // Assertion met
	}


	_, file, line, ok := runtime.Caller(skipFrames)

	// Could not get Caller info
	if !ok {
		panic(AssertionError{
			Message: msg,
		})
	}
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
