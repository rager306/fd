## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-14 - Auth Timing Attack Vulnerability
**Vulnerability:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the given slices differ, allowing a length-based timing leak on the `FD_API_KEY`.
**Learning:** Checking equality of different length slices bypasses constant-time bounds and creates a timing oracle. This oracle could let attackers guess the API key length or structure over time.
**Prevention:** If the slice lengths are different, balance the execution time by performing a dummy execution `subtle.ConstantTimeCompare(apiKey, apiKey)` to match execution lengths before returning an error.
