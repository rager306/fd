## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Leaks in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` leaks the length of the string if the two strings have different lengths, which could allow attackers to guess the API key length.
**Learning:** Hashing the inputs (e.g., with SHA-256) guarantees fixed-length slices for comparison, eliminating the length-based timing leak.
**Prevention:** Always hash the secret and the provided token before calling `subtle.ConstantTimeCompare`. Hash static expected secrets once during initialization.
