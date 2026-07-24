## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-24 - Length-Based Timing Attack in ConstantTimeCompare
**Vulnerability:** Using `subtle.ConstantTimeCompare` directly on raw strings/bytes where lengths may differ.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices are different. This creates a timing oracle that allows an attacker to discover the exact length of the secret token.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to ensure they are of equal length.
