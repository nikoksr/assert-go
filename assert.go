package assert

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type AssertionError struct {
	Message       string
	File          string
	Line          int
	SourceContext string
}

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

// getSourceContext reads the source file and returns lines around the failure.
func getSourceContext(file string, line int, contextLines int) string {
	f, err := os.Open(file)
	if err != nil {
		return "" // Return empty if we can't read the file
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Store the lines we want
	start := max(1, line-contextLines)
	end := line + contextLines

	var lines []string
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= start && currentLine <= end {
			prefix := "  "
			if currentLine == line {
				prefix = "→ " // Arrow pointing to the failure line
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

func Assert(condition bool, msg string, values ...any) {
	if !condition {
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

		// Get source context (5 lines before and after)
		sourceContext := getSourceContext(file, line, 5)

		err := AssertionError{
			Message:       msg + dumpInfo,
			File:          filepath.Base(file),
			Line:          line,
			SourceContext: sourceContext,
		}

		panic(err)
	}
}
