package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fd-api/handlers"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

func TestLifecycleGateRejectsBeforeWarmup(t *testing.T) {
	state := lifecycle.NewState()
	w := serveLifecycleRequest(state)

	assertLifecycleError(t, w, http.StatusServiceUnavailable, handlers.CodeModelNotLoaded, warmupRetryAfterSeconds)
}

func TestLifecycleGateRejectsDuringShutdown(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	state.BeginShutdown()
	w := serveLifecycleRequest(state)

	assertLifecycleError(t, w, http.StatusServiceUnavailable, handlers.CodeShuttingDown, shutdownRetryAfterSeconds)
}

func TestLifecycleGateTracksAcceptedRequest(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	r.GET("/protected", LifecycleGate(state), func(c *gin.Context) {
		close(started)
		<-release
		c.Status(http.StatusNoContent)
	})

	go func() {
		defer close(done)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/protected", http.NoBody))
		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d; body=%s", w.Code, http.StatusNoContent, w.Body.String())
		}
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("handler did not start")
	}
	if err := state.WaitDrain(0); !errors.Is(err, lifecycle.ErrDrainTimeout) {
		t.Fatalf("WaitDrain while request active = %v, want %v", err, lifecycle.ErrDrainTimeout)
	}

	close(release)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not finish")
	}
	if err := state.WaitDrain(0); err != nil {
		t.Fatalf("WaitDrain after response = %v, want nil", err)
	}
}

func serveLifecycleRequest(state *lifecycle.State) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", LifecycleGate(state), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/protected", http.NoBody))
	return w
}

func assertLifecycleError(t *testing.T, w *httptest.ResponseRecorder, status int, code, retryAfter string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, status, w.Body.String())
	}
	if got := w.Header().Get("Retry-After"); got != retryAfter {
		t.Fatalf("Retry-After = %q, want %q", got, retryAfter)
	}

	var body handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if body.Error.Code != code {
		t.Fatalf("error.code = %q, want %q", body.Error.Code, code)
	}
	if body.Error.Type != handlers.TypeOverloadedError {
		t.Fatalf("error.type = %q, want %q", body.Error.Type, handlers.TypeOverloadedError)
	}
}
