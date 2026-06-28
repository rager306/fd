## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-24 - Fix Length-Based Timing Attack in API Key Validation
**Vulnerability:** API key validation was vulnerable to a length-based timing attack because `subtle.ConstantTimeCompare` returns early if the inputs are different lengths.
**Learning:** When comparing secrets like API keys or tokens with `subtle.ConstantTimeCompare`, inputs must be guaranteed to have the same length to ensure execution time is truly constant.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) to ensure equal length before comparing them with `subtle.ConstantTimeCompare`.
