## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-20 - Length-based Timing Attack in API Key Validation
**Vulnerability:** Found `subtle.ConstantTimeCompare` comparing variable-length inputs without prior hashing in `api/middleware/auth.go`.
**Learning:** `ConstantTimeCompare` returns early if slice lengths differ, failing to provide constant-time guarantees and allowing length-based timing attacks.
**Prevention:** Always hash secrets using a cryptographically secure hash function (e.g., `crypto/sha256`) to ensure fixed-length inputs before passing to `ConstantTimeCompare`.
