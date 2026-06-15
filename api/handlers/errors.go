package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Canonical, machine-readable error codes (Section 3 of docs/fd-v2.md).
// Adding a new code here requires adding a row in errorCodeRegistry AND a test
// in errors_test.go.
const (
	CodeInputRequired     = "input_required"
	CodeInputTooLong      = "input_too_long"
	CodeBatchTooLarge     = "batch_too_large"
	CodeDimensionsInvalid = "dimensions_invalid"
	CodeInvalidJSON       = "invalid_json"
	CodeUnauthorized      = "unauthorized"
	CodeNotFound          = "not_found"
	CodeMethodNotAllowed  = "method_not_allowed"
	CodePayloadTooLarge   = "payload_too_large"
	CodeRateLimitExceeded = "rate_limit_exceeded"
	CodeInternalError     = "internal_error"
	CodeModelNotLoaded    = "model_not_loaded"
	CodeModelOverloaded   = "model_overloaded"
	CodeShuttingDown      = "shutting_down"
	CodeEncodingInvalid   = "encoding_format_invalid"
	CodePriorityInvalid   = "priority_invalid"
)

// OpenAI-style error types. Values match docs/fd-v2.md Section 3 catalog.
const (
	TypeInvalidRequest  = "invalid_request_error"
	TypeAuthError       = "authentication_error"
	TypePermissionError = "permission_error"
	TypeNotFoundError   = "not_found_error"
	TypeRateLimitError  = "rate_limit_error"
	TypeOverloadedError = "overloaded_error"
	TypeInternalError   = "internal_error"
)

// ErrorDetail is the inner envelope of an OpenAI-style error response.
type ErrorDetail struct {
	Code    string `json:"code"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse is the wire-level error envelope returned by fd v2.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// errorCodeRegistry maps canonical codes to (type, httpStatus).
// Single source of truth — WriteError consults this table only.
var errorCodeRegistry = map[string]struct {
	Type       string
	HTTPStatus int
}{
	CodeInputRequired:     {TypeInvalidRequest, http.StatusBadRequest},
	CodeInputTooLong:      {TypeInvalidRequest, http.StatusRequestEntityTooLarge},
	CodeBatchTooLarge:     {TypeInvalidRequest, http.StatusRequestEntityTooLarge},
	CodeDimensionsInvalid: {TypeInvalidRequest, http.StatusBadRequest},
	CodeEncodingInvalid:   {TypeInvalidRequest, http.StatusBadRequest},
	CodePriorityInvalid:   {TypeInvalidRequest, http.StatusBadRequest},
	CodeInvalidJSON:       {TypeInvalidRequest, http.StatusBadRequest},
	CodeUnauthorized:      {TypeAuthError, http.StatusUnauthorized},
	CodeNotFound:          {TypeNotFoundError, http.StatusNotFound},
	CodeMethodNotAllowed:  {TypeInvalidRequest, http.StatusMethodNotAllowed},
	CodePayloadTooLarge:   {TypeInvalidRequest, http.StatusRequestEntityTooLarge},
	CodeRateLimitExceeded: {TypeRateLimitError, http.StatusTooManyRequests},
	CodeInternalError:     {TypeInternalError, http.StatusInternalServerError},
	CodeModelNotLoaded:    {TypeOverloadedError, http.StatusServiceUnavailable},
	CodeModelOverloaded:   {TypeOverloadedError, http.StatusServiceUnavailable},
	CodeShuttingDown:      {TypeOverloadedError, http.StatusServiceUnavailable},
}

// HTTPStatusFor returns the canonical HTTP status for a code, or 500 if unknown.
func HTTPStatusFor(code string) int {
	if entry, ok := errorCodeRegistry[code]; ok {
		return entry.HTTPStatus
	}
	return http.StatusInternalServerError
}

// WriteError emits an OpenAI-style error envelope and aborts the gin context.
// param is the field name that triggered the error (e.g. "input", "dimensions");
// pass empty string when no specific param applies.
func WriteError(c *gin.Context, code, param, message string) {
	entry, ok := errorCodeRegistry[code]
	if !ok {
		// Unknown code — fail closed as internal_error so we never emit a
		// non-canonical envelope.
		entry = errorCodeRegistry[CodeInternalError]
		code = CodeInternalError
	}
	c.Set(ContextKeyErrorCode, code)
	c.AbortWithStatusJSON(entry.HTTPStatus, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Type:    entry.Type,
			Param:   param,
			Message: message,
		},
	})
}

// WriteErrorWithRetryAfter is like WriteError but also sets the Retry-After
// header (RFC 7231 §7.1.3). Used for 429/503 responses per R-P0-16.
func WriteErrorWithRetryAfter(c *gin.Context, code, param, message, retryAfter string) {
	c.Header("Retry-After", retryAfter)
	WriteError(c, code, param, message)
}

// AllErrorCodes returns the registered codes in deterministic order.
// Used by tests to assert coverage.
func AllErrorCodes() []string {
	return []string{
		CodeInputRequired,
		CodeInputTooLong,
		CodeBatchTooLarge,
		CodeDimensionsInvalid,
		CodeEncodingInvalid,
		CodePriorityInvalid,
		CodeInvalidJSON,
		CodeUnauthorized,
		CodeNotFound,
		CodeMethodNotAllowed,
		CodePayloadTooLarge,
		CodeRateLimitExceeded,
		CodeInternalError,
		CodeModelNotLoaded,
		CodeModelOverloaded,
		CodeShuttingDown,
	}
}

const (
	// ContextKeyValidatedRequest is the gin context key the validation
	// middleware uses to pass the parsed embed.EmbeddingsRequest downstream.
	// Lives in handlers (not middleware) to avoid a middleware→handlers import
	// cycle from the /v1/embeddings handler that needs to look it up.
	ContextKeyValidatedRequest = "fd_validated_embeddings_request"
	// ContextKeyErrorCode stores the canonical fd error code emitted by
	// WriteError so observability middleware can increment fd_errors_total.
	ContextKeyErrorCode = "fd_error_code"
)
