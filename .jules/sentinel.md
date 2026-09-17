## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-28 - Timing Leak in API Key Comparison
**Vulnerability:** Length-based timing leak in API key authentication because `subtle.ConstantTimeCompare` returns immediately if slice lengths differ.
**Learning:** When comparing secrets of potentially arbitrary or varying lengths using `subtle.ConstantTimeCompare`, both values must be hashed first to guarantee fixed-length slices.
**Prevention:** Hash static expected secrets once during middleware initialization and hash incoming tokens on each request before comparing them.
