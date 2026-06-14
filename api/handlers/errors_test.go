package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test fixtures — param names are repeated across table-driven cases
// below. Constants prevent goconst from flagging these as duplicates while
// keeping the test readable.
const (
	paramInput          = "input"
	paramDimensions     = "dimensions"
	paramEncodingFormat = "encoding_format"
	paramPriority       = "priority"
)

func TestAllErrorCodesRegistered(t *testing.T) {
	// Every code returned by AllErrorCodes must exist in the registry.
	for _, code := range AllErrorCodes() {
		if _, ok := errorCodeRegistry[code]; !ok {
			t.Errorf("code %q listed in AllErrorCodes but missing from registry", code)
		}
	}
}

func TestErrorEnvelopeShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		code     string
		param    string
		message  string
		wantCode string
		wantType string
		wantHTTP int
	}{
		{CodeInputRequired, paramInput, "input is required", "input_required", TypeInvalidRequest, http.StatusBadRequest},
		{CodeInputTooLong, paramInput, "input[0] exceeds max length 512 tokens", "input_too_long", TypeInvalidRequest, http.StatusRequestEntityTooLarge},
		{CodeBatchTooLarge, paramInput, "batch size 100 exceeds max 32", "batch_too_large", TypeInvalidRequest, http.StatusRequestEntityTooLarge},
		{CodeDimensionsInvalid, paramDimensions, "dimensions must be 1024 or 512, got 99999", "dimensions_invalid", TypeInvalidRequest, http.StatusBadRequest},
		{CodeDimensionsRequired, paramDimensions, "dimensions is required", "dimensions_required", TypeInvalidRequest, http.StatusBadRequest},
		{CodeDimensionsMismatch, paramDimensions, "model does not support 512-dim", "dimensions_mismatch", TypeInvalidRequest, http.StatusBadRequest},
		{CodeEncodingInvalid, paramEncodingFormat, "encoding_format must be float or base64", "encoding_format_invalid", TypeInvalidRequest, http.StatusBadRequest},
		{CodePriorityInvalid, paramPriority, "priority must be low, normal, or high", "priority_invalid", TypeInvalidRequest, http.StatusBadRequest},
		{CodeInvalidJSON, "", "invalid JSON: unexpected end of JSON input", "invalid_json", TypeInvalidRequest, http.StatusBadRequest},
		{CodeUnauthorized, "", "missing or invalid API key", "unauthorized", TypeAuthError, http.StatusUnauthorized},
		{CodeNotFound, "", "path /v9999 not found", "not_found", TypeNotFoundError, http.StatusNotFound},
		{CodePayloadTooLarge, "", "request body 52428800 bytes exceeds max 10485760 bytes", "payload_too_large", TypeInvalidRequest, http.StatusRequestEntityTooLarge},
		{CodeRateLimitExceeded, "", "rate limit exceeded; retry after 60s", "rate_limit_exceeded", TypeRateLimitError, http.StatusTooManyRequests},
		{CodeInternalError, "", "internal server error; request_id=abc", "internal_error", TypeInternalError, http.StatusInternalServerError},
		{CodeModelNotLoaded, "", "model not loaded; retry after 5s", "model_not_loaded", TypeOverloadedError, http.StatusServiceUnavailable},
		{CodeModelOverloaded, "", "model overloaded; retry after 5s", "model_overloaded", TypeOverloadedError, http.StatusServiceUnavailable},
		{CodeShuttingDown, "", "service shutting down; retry after 30s", "shutting_down", TypeOverloadedError, http.StatusServiceUnavailable},
		{CodeRequestTimeout, "", "request timed out after 30s", "request_timeout", TypeOverloadedError, http.StatusGatewayTimeout},
	}

	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			WriteError(c, tc.code, tc.param, tc.message)

			if w.Code != tc.wantHTTP {
				t.Errorf("HTTP status = %d, want %d", w.Code, tc.wantHTTP)
			}
			var resp ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal envelope: %v body=%s", err, w.Body.String())
			}
			if resp.Error.Code != tc.wantCode {
				t.Errorf("code = %q, want %q", resp.Error.Code, tc.wantCode)
			}
			if resp.Error.Type != tc.wantType {
				t.Errorf("type = %q, want %q", resp.Error.Type, tc.wantType)
			}
			if resp.Error.Param != tc.param {
				t.Errorf("param = %q, want %q", resp.Error.Param, tc.param)
			}
			if resp.Error.Message != tc.message {
				t.Errorf("message = %q, want %q", resp.Error.Message, tc.message)
			}
		})
	}
}

func TestWriteErrorUnknownCodeFailsClosedAsInternal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	WriteError(c, "bogus_code_we_never_registered", "input", "test")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("unknown code HTTP = %d, want 500", w.Code)
	}
	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error.Code != CodeInternalError {
		t.Errorf("unknown code mapped to %q, want %q", resp.Error.Code, CodeInternalError)
	}
	if resp.Error.Type != TypeInternalError {
		t.Errorf("unknown code type = %q, want %q", resp.Error.Type, TypeInternalError)
	}
}

func TestWriteErrorWithRetryAfter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	WriteErrorWithRetryAfter(c, CodeShuttingDown, "", "service shutting down", "30")

	if got := w.Header().Get("Retry-After"); got != "30" {
		t.Errorf("Retry-After = %q, want 30", got)
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("HTTP = %d, want 503", w.Code)
	}
	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error.Code != CodeShuttingDown {
		t.Errorf("code = %q, want %q", resp.Error.Code, CodeShuttingDown)
	}
}

func TestWriteErrorAbortsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	WriteError(c, CodeInputRequired, "input", "missing")

	if !c.IsAborted() {
		t.Error("expected context to be aborted after WriteError")
	}
}

func TestHTTPStatusForUnknownReturns500(t *testing.T) {
	if got := HTTPStatusFor("not_a_real_code"); got != http.StatusInternalServerError {
		t.Errorf("HTTPStatusFor(unknown) = %d, want 500", got)
	}
}
