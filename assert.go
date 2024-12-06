package assert

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Config is used to configure the behavior of the assertion library.
type Config struct {
	// IncludeSource determines if source context is included in errors.
	//
	// Default: true
	IncludeSource bool

	// ContextLines determines how many lines of context to show.
	//
	// Default: 5
	ContextLines int
}

//nolint:gochecknoglobals // read-only global
var activeConfig = Config{
	IncludeSource: true,
	ContextLines:  5, //nolint:gomnd,mnd // default value
}

// SetConfig sets the configuration for the assertion library.
func SetConfig(config Config) {
	activeConfig = config
}

// Assert is the main assertion function. It panics if the condition is false.
type AssertionError struct {
	Message       string
	File          string
	SourceContext string
	Line          int
}

// Error returns the error message.
func (e AssertionError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Assertion failed at %s:%d\n", e.File, e.Line))
	sb.WriteString(fmt.Sprintf("Message: %s\n", e.Message))

	if e.SourceContext != "" {
		sb.WriteString("Source context:\n")
		sb.WriteString(e.SourceContext)
	}

	return sb.String()
}

// Assert panics if the condition is false. Configurable via SetConfig.
func Assert(condition bool, msg string, values ...any) {
	if condition {
		return // Assertion met
	}

	_, file, line, _ := runtime.Caller(1)

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
