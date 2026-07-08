## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-18 - Timing attack via length checks in ConstantTimeCompare
**Vulnerability:** Comparing plaintext secrets directly with `subtle.ConstantTimeCompare` allows length-based timing attacks, because the function checks length first and returns early if they differ.
**Learning:** To guarantee true constant-time comparison, both inputs must first be hashed to ensure they are the exact same length.
**Prevention:** Always hash secrets (e.g. using `crypto/sha256.Sum256`) before passing them to `subtle.ConstantTimeCompare`.
