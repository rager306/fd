## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-23 - Timing Attack via ConstantTimeCompare
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the provided token and the expected API key. `subtle.ConstantTimeCompare` returns immediately if the lengths of the slices differ, exposing a timing side-channel that leak the exact length of the expected API key.
**Learning:** In Go, variable-length secrets (like API tokens) must not be compared directly with `ConstantTimeCompare`.
**Prevention:** Always hash both the expected secret and the provided input (e.g., using `crypto/sha256.Sum256`) and compare the resulting fixed-length hashes to ensure comparison time is independent of the input length.
