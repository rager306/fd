package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const paramInput = "input"

func TestLimitRequestBodyRejectsDeclaredOversizeBeforeDownstream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	downstreamCalled := false
	r := gin.New()
	r.POST("/batch", LimitRequestBody(), func(c *gin.Context) {
		downstreamCalled = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/batch", strings.NewReader(`{"inputs":["ok"]}`))
	req.ContentLength = maxRequestBodyBytes + 1
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertCode(t, w, handlers.CodePayloadTooLarge, http.StatusRequestEntityTooLarge, "")
	if downstreamCalled {
		t.Fatal("downstream handler should not be called for oversized request")
	}
}

// runMiddleware runs a single request through the validation middleware
// and returns the recorder. The middleware is wired to a downstream
// handler that simply stores the validated request in the context for
// later inspection (so we can assert both the error path AND the happy
// path delivers the parsed request to downstream handlers).
func runMiddleware(t *testing.T, body string, contentLen int64) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", bodyReader)
	if contentLen > 0 {
		req.ContentLength = contentLen
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req // gin v1.12 CreateTestContext does NOT set c.Request
	ValidateEmbeddingsRequest()(c)
	// If middleware did NOT abort, simulate the downstream handler.
	if !c.IsAborted() {
		validated, ok := c.Get(handlers.ContextKeyValidatedRequest)
		if !ok {
			t.Fatal("middleware passed but ContextKeyValidatedRequest not set")
		}
		_ = validated.(*embed.EmbeddingsRequest)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
	return w
}

func TestValidationHappyPathFloat(t *testing.T) {
	body := `{"input":["hello","world"],"dimensions":1024}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationHappyPathBase64(t *testing.T) {
	body := `{"input":["hello"],"encoding_format":"base64"}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationHappyPathPriorityAndUser(t *testing.T) {
	body := `{"input":["hello"],"priority":"high","user":"caller-123"}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationEmptyInput(t *testing.T) {
	body := `{"input":[]}`
	w := runMiddleware(t, body, int64(len(body)))
	assertCode(t, w, handlers.CodeInputRequired, http.StatusBadRequest, "input")
}

func TestValidationMissingInput(t *testing.T) {
	// {} — input field absent
	w := runMiddleware(t, `{}`, 2)
	assertCode(t, w, handlers.CodeInputRequired, http.StatusBadRequest, "input")
}

func TestValidationBatchTooLarge(t *testing.T) {
	// 129 > maxBatchSize=128 (raised June 2026 to better match OpenAI v1 behavior;
	// TEI sub-batches of 32 are handled in the handler).
	inputs := make([]string, 129)
	for i := range inputs {
		inputs[i] = "x"
	}
	body, _ := json.Marshal(gin.H{paramInput: inputs})
	w := runMiddleware(t, string(body), int64(len(body)))
	assertCode(t, w, handlers.CodeBatchTooLarge, http.StatusRequestEntityTooLarge, paramInput)
}

func TestValidationBatchAt128Boundary(t *testing.T) {
	// 128 == maxBatchSize, should pass validation (TEI chunking handles).
	inputs := make([]string, 128)
	for i := range inputs {
		inputs[i] = "x"
	}
	body, _ := json.Marshal(gin.H{"input": inputs})
	w := runMiddleware(t, string(body), int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("batch=128 should pass; got status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestValidationInputTooLong(t *testing.T) {
	// 2049 chars > maxInputChars=2048
	longText := strings.Repeat("x", 2049)
	body, _ := json.Marshal(gin.H{"input": []string{longText}})
	w := runMiddleware(t, string(body), int64(len(body)))
	assertCode(t, w, handlers.CodeInputTooLong, http.StatusRequestEntityTooLarge, "input")
}

func TestValidationDimensionsInvalid(t *testing.T) {
	body := `{"input":["x"],"dimensions":99999}`
	w := runMiddleware(t, body, int64(len(body)))
	assertCode(t, w, handlers.CodeDimensionsInvalid, http.StatusBadRequest, "dimensions")
}

func TestValidationDimensions512(t *testing.T) {
	body := `{"input":["x"],"dimensions":512}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationDimensions1024(t *testing.T) {
	body := `{"input":["x"],"dimensions":1024}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationDimensionsMissing(t *testing.T) {
	body := `{"input":["x"]}` // no dimensions → default 1024 in handler
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestValidationEncodingFormatFloat(t *testing.T) {
	body := `{"input":["x"],"encoding_format":"float"}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestValidationEncodingFormatBase64(t *testing.T) {
	body := `{"input":["x"],"encoding_format":"base64"}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestValidationEncodingFormatInvalid(t *testing.T) {
	body := `{"input":["x"],"encoding_format":"hex"}`
	w := runMiddleware(t, body, int64(len(body)))
	assertCode(t, w, handlers.CodeEncodingInvalid, http.StatusBadRequest, "encoding_format")
}

func TestValidationPriorityInvalid(t *testing.T) {
	body := `{"input":["x"],"priority":"urgent"}`
	w := runMiddleware(t, body, int64(len(body)))
	assertCode(t, w, handlers.CodePriorityInvalid, http.StatusBadRequest, "priority")
}

func TestValidationInvalidJSONMalformed(t *testing.T) {
	body := `{bad json`
	w := runMiddleware(t, body, int64(len(body)))
	assertCode(t, w, handlers.CodeInvalidJSON, http.StatusBadRequest, "")
}

func TestValidationInvalidJSONNonString(t *testing.T) {
	// input:[123] — array element must be string. Per spec T-E-4 this must
	// return 400 invalid_request_error, NOT the leaky "json: cannot unmarshal"
	// from the original embedding handler. Our middleware emits a clean
	// invalid_json envelope because the JSON itself parses as a valid
	// []interface{} but fails type-check inside the request struct.
	//
	// The EmbeddingsRequest.UnmarshalJSON accepts only string or []string
	// for input. A []interface{} like [123] would be caught at the JSON
	// decoder level before our type-specific check. So we get invalid_json
	// here. That's the right behavior: callers should not see Go internals.
	body := `{"input":[123]}`
	w := runMiddleware(t, body, int64(len(body)))
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal envelope: %v; body=%s", err, w.Body.String())
	}
	if resp.Error.Message != "input must be an array of strings, got array" {
		t.Fatalf("message = %q, want clean array element message", resp.Error.Message)
	}
	if strings.Contains(resp.Error.Message, "input[") {
		t.Fatalf("message should not contain malformed input[ prefix: %q", resp.Error.Message)
	}

	// Body must NOT contain the leaky "json: cannot unmarshal" string.
	if strings.Contains(w.Body.String(), "json: cannot unmarshal") {
		t.Errorf("body leaks Go internals: %s", w.Body.String())
	}
}

func TestValidationPayloadTooLarge(t *testing.T) {
	// Construct a body whose declared Content-Length exceeds the 10MB cap.
	// http.MaxBytesReader rejects reads past the limit; the subsequent
	// ShouldBindJSON returns *http.MaxBytesError which we map to
	// payload_too_large. We pass an explicitly huge Content-Length and a
	// truncated reader so the error is deterministic.
	big := bytes.Repeat([]byte("x"), 100) // actual body is small
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", bytes.NewReader(big))
	req.ContentLength = int64(maxRequestBodyBytes + 1) // declared as oversized
	req.Header.Set("Content-Type", "application/json")
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	ValidateEmbeddingsRequest()(c)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413; body=%s", w.Code, w.Body.String())
	}
}

// assertCode verifies the response carries the expected OpenAI error envelope.
func assertCode(t *testing.T, w *httptest.ResponseRecorder, wantCode string, wantStatus int, wantParam string) {
	t.Helper()
	if w.Code != wantStatus {
		t.Errorf("HTTP = %d, want %d; body=%s", w.Code, wantStatus, w.Body.String())
		return
	}
	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal envelope: %v; body=%s", err, w.Body.String())
	}
	if resp.Error.Code != wantCode {
		t.Errorf("error.code = %q, want %q", resp.Error.Code, wantCode)
	}
	if resp.Error.Param != wantParam {
		t.Errorf("error.param = %q, want %q", resp.Error.Param, wantParam)
	}
	// Sanity: type must be invalid_request_error for 4xx envelopes
	if wantStatus >= 400 && wantStatus < 500 && resp.Error.Type != handlers.TypeInvalidRequest {
		t.Errorf("error.type = %q, want %q", resp.Error.Type, handlers.TypeInvalidRequest)
	}
}
