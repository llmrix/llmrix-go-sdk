// Package apierr defines the SDK's error hierarchy.
// These types are re-exported from the root llmrix package via type aliases;
// callers should reference them as llmrix.LlmrixError, llmrix.LlmrixApiError, etc.
package apierr

import "fmt"

// LlmrixError is the base error type returned by all SDK operations.
type LlmrixError struct {
	// Message is a human-readable summary of what went wrong.
	Message string
	// Cause holds the underlying error, if any.
	Cause error
}

func (e *LlmrixError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("llmrix: %s: %v", e.Message, e.Cause)
	}
	return "llmrix: " + e.Message
}

func (e *LlmrixError) Unwrap() error { return e.Cause }

// LlmrixApiError is returned when the server responds with an HTTP 4xx/5xx status.
type LlmrixApiError struct {
	LlmrixError
	// StatusCode is the HTTP response status code.
	StatusCode int
	// ResponseBody is the raw response body; may be empty.
	ResponseBody string
}

func (e *LlmrixApiError) Error() string {
	return fmt.Sprintf("llmrix: HTTP %d: %s", e.StatusCode, e.Message)
}

// LlmrixAuthError is returned on HTTP 401 (invalid API key) or 503
// (server running in native mode without auth configured).
type LlmrixAuthError struct {
	LlmrixApiError
}

// New creates a plain LlmrixError.
func New(msg string, cause error) error {
	return &LlmrixError{Message: msg, Cause: cause}
}

// NewAPI creates a LlmrixApiError, promoting to LlmrixAuthError for 401/503.
func NewAPI(code int, msg, body string) error {
	base := LlmrixApiError{
		LlmrixError:  LlmrixError{Message: msg},
		StatusCode:   code,
		ResponseBody: body,
	}
	if code == 401 || code == 503 {
		return &LlmrixAuthError{base}
	}
	return &base
}
