## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - Timing Attack via ConstantTimeCompare
**Vulnerability:** API key comparison using `subtle.ConstantTimeCompare` directly on byte slices of potentially different lengths.
**Learning:** `ConstantTimeCompare` returns immediately if lengths differ, leaking the secret's length.
**Prevention:** Hash both inputs with a fixed-length hash (like SHA-256) before using `ConstantTimeCompare`.
