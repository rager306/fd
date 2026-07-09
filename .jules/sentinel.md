## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-09 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** Used `subtle.ConstantTimeCompare` directly on tokens of potentially different lengths.
**Learning:** `ConstantTimeCompare` returns early if lengths differ, making it vulnerable to length-based timing attacks.
**Prevention:** Hash both inputs (e.g., using `crypto/sha256`) to ensure equal length before comparing them with `subtle.ConstantTimeCompare`.
