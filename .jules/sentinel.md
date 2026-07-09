## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-27 - Authentication Timing Attack
**Vulnerability:** The API key validation in `api/middleware/auth.go` uses `subtle.ConstantTimeCompare` directly on the provided token and the configured API key strings.
**Learning:** `subtle.ConstantTimeCompare` only protects against timing attacks if both byte slices are the same length. If lengths differ, it returns early.
**Prevention:** To guarantee equal length and prevent length-based timing attacks, both inputs should be hashed (e.g., using `crypto/sha256`) before comparison.

## 2024-05-27 - Fixing Lint Errors on Existing Fix
**Vulnerability:** Not a vulnerability, but fixing linter issues regarding `errcheck` on previous code changes in CI.
**Learning:** CI failures on previous commits are blocking the PR check. When fixing `errcheck` in CI, ensure all deferred `Close()` calls have the error ignored using `_ = file.Close()`.
**Prevention:** Verify `golangci-lint` passes on CI before merging.
