## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-18 - Length-based Timing Attack in Token Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on variable-length API tokens and the expected API key, which leaks the exact length of the expected API key via early return if lengths differ.
**Learning:** `subtle.ConstantTimeCompare` does not protect against timing attacks if the lengths of the two byte slices are different (it returns immediately).
**Prevention:** Always hash both the expected secret and the user-provided token (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to guarantee they are of equal length.
