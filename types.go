package assert

import (
	"fmt"
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

// AssertionError is the error type returned when an assertion fails.
type AssertionError struct {
	Message       string
	File          string
	SourceContext string
	Line          int
}

// Error returns the error message.
func (e AssertionError) Error() string {
	var sb strings.Builder

	if e.File != "" {
		sb.WriteString(fmt.Sprintf("Assertion failed at %s:%d\n", e.File, e.Line))
	} else {
		sb.WriteString("Assertion failed (Runtime caller info is not available)\n")
	}
	sb.WriteString(fmt.Sprintf("Message: %s\n", e.Message))

	if e.SourceContext != "" {
		sb.WriteString("Source context:\n")
		sb.WriteString(e.SourceContext)
	}

	return sb.String()
}
