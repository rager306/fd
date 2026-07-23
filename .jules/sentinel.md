## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Attack in Comparison
**Vulnerability:** API key verification used `subtle.ConstantTimeCompare` directly on raw slices without hashing.
**Learning:** `subtle.ConstantTimeCompare` leaks length information by returning early when lengths differ.
**Prevention:** Always hash both inputs (e.g., `crypto/sha256`) to ensure equal lengths before calling `ConstantTimeCompare`.
