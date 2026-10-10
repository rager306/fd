## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-03-24 - Timing Oracle in token comparison
**Vulnerability:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, creating a length-based timing oracle when validating bearer tokens.
**Learning:** `subtle.ConstantTimeCompare` is only constant time for slices of the same length. Comparing secrets of different lengths against user input leaks length information.
**Prevention:** Mitigate the timing leak by executing a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths mismatch to strictly balance the execution time.
