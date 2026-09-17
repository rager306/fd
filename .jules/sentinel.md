## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-17 - Timing Attack in Auth Middleware
**Vulnerability:** Length-based timing leak in API key comparison.
**Learning:** Using `subtle.ConstantTimeCompare` on strings of different lengths allows an attacker to deduce the length of the expected secret.
**Prevention:** Hash both the expected secret and the incoming token (e.g., using SHA-256) before comparison to ensure fixed-length slices.
