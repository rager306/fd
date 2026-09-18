## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-18 - Timing Leak in API Key Auth
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the provided token and expected API key strings.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices differ. This leaks the exact length of the expected API key via timing attacks.
**Prevention:** To prevent length-based timing leaks, securely hash both the expected secret and the incoming token (e.g., using SHA-256) before passing them to `ConstantTimeCompare`. Compute the expected hash once during initialization for performance.
