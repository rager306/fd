## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-26 - Timing Attack Vulnerability in API Key Comparison
**Vulnerability:** The API key validation in `api/middleware/auth.go` used `subtle.ConstantTimeCompare` directly on the provided token and the configured API key. This is vulnerable to length-based timing attacks because `ConstantTimeCompare` returns immediately if the lengths of the two byte slices differ, leaking the length of the secret API key.
**Learning:** Even when using "constant time" functions, one must be aware of their preconditions. `subtle.ConstantTimeCompare` only provides constant-time properties for inputs of equal length.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare`. This guarantees they have the same length and completely mitigates length-based timing attacks.
