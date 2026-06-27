## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-27 - Length-based Timing Attack in Token Comparison
**Vulnerability:** Comparing raw tokens using `subtle.ConstantTimeCompare` without hashing.
**Learning:** `ConstantTimeCompare` returns early if the lengths of the two byte slices differ, which can lead to length-based timing attacks. By hashing both inputs first, we ensure they are of equal length, allowing `ConstantTimeCompare` to execute in constant time regardless of the original token lengths.
**Prevention:** When comparing secrets or tokens using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length.
