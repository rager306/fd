## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-24 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** The APIKeyAuth middleware directly compared the incoming bearer token with the expected API key using `subtle.ConstantTimeCompare`.
**Learning:** `subtle.ConstantTimeCompare` only provides constant time execution when the input lengths match. If lengths differ, it instantly returns, revealing length information.
**Prevention:** Hash both the expected secret and the incoming token (e.g., using SHA-256) before comparison to guarantee fixed-length slices and prevent length-based leaks.
