// Package middleware contains gin middleware for fd v2 request validation
// and observability. Middleware run in this order (outermost first):
//
//  1. headers (S03)        — X-Request-Id, Server, etc
//  2. metrics (S03)         — request counters/histograms
//  3. validation (T03 here) — body size, JSON shape, input/dimensions rules
//  4. auth (S05)            — FD_API_KEY bearer (skips /live, /metrics, /docs)
//  5. rate-limit (S05)      — per-IP/per-user
//  6. lifecycle (S02)       — warmup/shutdown gates
//  7. cache (S04)           — X-Cache HIT/MISS
//
// validation must run BEFORE model/cache lookups so 4xx/5xx are emitted
// without burning inference capacity (root cause of B9: 100 inputs → 500).
package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const (
	// maxRequestBodyBytes is the upper bound for the request body. Matches
	// the 10MB cap from docs/fd-v2.md Section 3 (payload_too_large).
	maxRequestBodyBytes = 10 * 1024 * 1024

	// maxBatchSize is the maximum number of inputs per /v1/embeddings call.
	// Updated to 128 (June 2026) to better match OpenAI v1 default behavior
	// (max 2048 per request per community discussion) while staying well
	// below TEI's max_client_batch_size=32 sub-limit. The handler chunks
	// internally into TEI-friendly sub-batches of 32.
	maxBatchSize = 128

	// Note: teiSubBatchSize (TEI's --max-client-batch-size=32) is defined
	// as a local const inside api/handlers/embeddings.go::CreateEmbedding
	// where the chunked loop runs. We don't export it from this package
	// because the middleware only validates input size; it does not
	// perform the actual chunking.

	// maxInputChars approximates 512 tokens at ~4 chars/token. Inputs longer
	// than this are rejected with 413 input_too_long instead of being
	// silently truncated by the model.
	maxInputChars = 2048

	// maxTotalInputChars caps the SUM of all input lengths in a single
	// request, mirroring OpenAI's 300,000 total-token-per-request limit
	// (4 chars/token ≈ 75000 chars cap, rounded up to be slightly more
	// permissive for callers using tokenizers that average <4 chars/token).
	maxTotalInputChars = 300_000 * 4
)

// ValidateEmbeddingsRequest is a gin middleware that:
//   - caps request body size at maxRequestBodyBytes (413 payload_too_large)
//   - parses JSON into embed.EmbeddingsRequest
//   - validates input array (non-empty, len<=32, all strings, each <=2048 chars)
//   - validates dimensions (must be 512 or 1024 if set)
//   - validates encoding_format (float|base64 if set) — T04 responsibility
//
// On any failure, the request is aborted with an OpenAI-style envelope
// (handlers.WriteError) so callers can switch on error.code.
func ValidateEmbeddingsRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Body size cap. http.MaxBytesReader returns *http.MaxBytesError
		// from any subsequent Read once the cap is exceeded; we surface
		// that specifically as 413 payload_too_large.
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodyBytes)

		if !validateContentLength(c) {
			return
		}
		req, ok := bindEmbeddingsRequest(c)
		if !ok {
			return
		}
		if !validateEmbeddingsPayload(c, req) {
			return
		}

		c.Set(handlers.ContextKeyValidatedRequest, req)
		c.Next()
	}
}

func validateContentLength(c *gin.Context) bool {
	// Reject early on declared Content-Length so we don't waste cycles
	// reading 10MB+ bodies just to throw them away. MaxBytesReader alone
	// only checks at read time, so a client declaring 50MB but sending 1KB
	// would slip through. This guards both honest and adversarial cases.
	if c.Request.ContentLength <= maxRequestBodyBytes {
		return true
	}
	handlers.WriteError(c, handlers.CodePayloadTooLarge, "",
		fmt.Sprintf("request body %d bytes exceeds max %d bytes", c.Request.ContentLength, maxRequestBodyBytes))
	return false
}

func bindEmbeddingsRequest(c *gin.Context) (*embed.EmbeddingsRequest, bool) {
	var req embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleBindError(c, err)
		return nil, false
	}
	return &req, true
}

func handleBindError(c *gin.Context, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		handlers.WriteError(c, handlers.CodePayloadTooLarge, "",
			fmt.Sprintf("request body exceeds max %d bytes", maxRequestBodyBytes))
		return
	}
	// Type mismatch (e.g. input:[123]) returns json.UnmarshalTypeError.
	// Map to a clean invalid_request_error envelope WITHOUT leaking
	// Go internals like "cannot unmarshal array into Go value".
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		handlers.WriteError(c, handlers.CodeInputRequired, "input",
			fmt.Sprintf("input[%s] must be string, got %s", typeErr.Field, typeErr.Value))
		return
	}
	handlers.WriteError(c, handlers.CodeInvalidJSON, "", "invalid JSON: "+err.Error())
}

func validateEmbeddingsPayload(c *gin.Context, req *embed.EmbeddingsRequest) bool {
	if len(req.Input) == 0 {
		handlers.WriteError(c, handlers.CodeInputRequired, "input",
			"input is required (non-empty array of strings)")
		return false
	}
	if len(req.Input) > maxBatchSize {
		handlers.WriteError(c, handlers.CodeBatchTooLarge, "input",
			fmt.Sprintf("batch size %d exceeds max %d; split into smaller batches", len(req.Input), maxBatchSize))
		return false
	}
	return validateInputLengths(c, req.Input) &&
		validateDimensions(c, req.Dimensions) &&
		validateEncodingFormat(c, req.EncodingFormat) &&
		validatePriority(c, req.Priority)
}

func validateInputLengths(c *gin.Context, input []string) bool {
	// Per-input length AND total length (mirrors OpenAI 300k token cap).
	totalChars := 0
	for i, t := range input {
		if len(t) > maxInputChars {
			handlers.WriteError(c, handlers.CodeInputTooLong, "input",
				fmt.Sprintf("input[%d] exceeds max length %d chars (got %d)", i, maxInputChars, len(t)))
			return false
		}
		totalChars += len(t)
	}
	if totalChars <= maxTotalInputChars {
		return true
	}
	handlers.WriteError(c, handlers.CodePayloadTooLarge, "input",
		fmt.Sprintf("total input length %d chars exceeds max %d chars (≈ OpenAI 300k token limit)", totalChars, maxTotalInputChars))
	return false
}

func validateDimensions(c *gin.Context, dimensions *int) bool {
	// dimensions: 512 or 1024 (or unset → default 1024 in handler).
	if dimensions == nil {
		return true
	}
	d := *dimensions
	if d == 512 || d == 1024 {
		return true
	}
	handlers.WriteError(c, handlers.CodeDimensionsInvalid, "dimensions",
		fmt.Sprintf("dimensions must be 1024 or 512, got %d", d))
	return false
}

func validateEncodingFormat(c *gin.Context, encodingFormat *string) bool {
	// encoding_format: float (default) or base64. Validated here so
	// T04 output shaping can rely on a small enum.
	if encodingFormat == nil {
		return true
	}
	ef := *encodingFormat
	if ef == "float" || ef == "base64" {
		return true
	}
	handlers.WriteError(c, handlers.CodeEncodingInvalid, "encoding_format",
		fmt.Sprintf("encoding_format must be float or base64, got %q", ef))
	return false
}

func validatePriority(c *gin.Context, priority *string) bool {
	// priority is accepted for OpenAI-compatible clients and future routing.
	// It is not used for scheduling yet, but validating the enum prevents
	// callers from persisting arbitrary values into future observability/rate-limit paths.
	if priority == nil || *priority == "" {
		return true
	}
	p := *priority
	if p == "low" || p == "normal" || p == "high" {
		return true
	}
	handlers.WriteError(c, handlers.CodePriorityInvalid, "priority",
		fmt.Sprintf("priority must be low, normal, or high, got %q", p))
	return false
}
