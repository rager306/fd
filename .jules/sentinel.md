## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-02 - Variable-Length Secret Timing Attack
**Vulnerability:** Go's `subtle.ConstantTimeCompare` returns immediately if byte slices have different lengths, exposing a timing side-channel that leaks the length of the secret API key.
**Learning:** Securely comparing variable-length secrets requires hashing both inputs to fixed-length hashes before comparison.
**Prevention:** Always hash secrets of potentially variable lengths (like API tokens) with SHA-256 before using `subtle.ConstantTimeCompare`, and pre-compute the expected secret's hash during middleware initialization to avoid per-request overhead.
