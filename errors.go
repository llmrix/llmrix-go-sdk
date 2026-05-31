package llmrix

import "github.com/llmrix/llmrix-go-sdk/internal/apierr"

// Re-export error types from internal/apierr so callers only need to import
// this package. errors.As / errors.Is work correctly with these aliases.

// LlmrixError is the base error type returned by all SDK operations.
// Use errors.As to check for more specific sub-types.
type LlmrixError = apierr.LlmrixError

// LlmrixApiError is returned when the server responds with an HTTP 4xx/5xx status.
// Inspect StatusCode to branch on specific conditions (e.g. 404, 429).
type LlmrixApiError = apierr.LlmrixApiError

// LlmrixAuthError is returned on HTTP 401 (invalid API key) or 503
// (server running in native mode without auth configured).
type LlmrixAuthError = apierr.LlmrixAuthError
