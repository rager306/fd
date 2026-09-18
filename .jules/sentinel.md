## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing attack vulnerability in token comparison
**Vulnerability:** Length-based timing leak in API key authentication
**Learning:** subtle.ConstantTimeCompare leaks length information if the slices being compared are of different lengths. A variable-length token allows an attacker to deduce the expected token length.
**Prevention:** Always hash tokens (e.g. SHA-256) before comparison to ensure fixed-length slices, and precompute the expected secret's hash to avoid overhead.
