## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-06 - Length-Based Timing Attack in API Key Auth
**Vulnerability:** API key authentication was susceptible to length-based timing attacks. `subtle.ConstantTimeCompare` was used to compare the provided token and the configured API key directly.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ, leaking the expected secret's length.
**Prevention:** When comparing secrets using `subtle.ConstantTimeCompare`, always hash both inputs (e.g., using `crypto/sha256`) first to guarantee equal length before comparison.
