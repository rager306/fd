## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-07 - Mitigate API Key Timing Leak
**Vulnerability:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, creating a length-based timing oracle that can leak the length of the valid API key.
**Learning:** Comparing secrets against user input using `subtle.ConstantTimeCompare` without checking lengths can introduce timing leaks. Hashing inputs could prevent this but introduces a DoS vulnerability.
**Prevention:** Always execute a dummy constant-time comparison when lengths mismatch to balance the execution time strictly.
