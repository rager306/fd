## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-22 - Length-Based Timing Attack in Token Comparison
**Vulnerability:** Comparing bearer tokens against the API key using `subtle.ConstantTimeCompare` without hashing inputs allows length-based timing attacks because the function returns early if slice lengths differ.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time comparison if the input lengths are identical. If lengths differ, it is not constant time.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) to ensure equal lengths before comparing secrets or tokens using `subtle.ConstantTimeCompare`.
