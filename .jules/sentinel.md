## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-13 - Length-Based Timing Attack in Token Comparison
**Vulnerability:** Comparing API tokens directly using `subtle.ConstantTimeCompare` without hashing them first.
**Learning:** `ConstantTimeCompare` checks length first and returns early if they don't match, exposing a length-based timing attack.
**Prevention:** Ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length before calling `ConstantTimeCompare`.
