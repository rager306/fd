## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-10-08 - Timing Attack Oracle in Authentication Middleware
**Vulnerability:** The authentication middleware was vulnerable to a length-based timing attack because `subtle.ConstantTimeCompare` returns immediately if slice lengths differ.
**Learning:** When comparing secrets against user input using `subtle.ConstantTimeCompare`, differing lengths create a timing oracle.
**Prevention:** Mitigate the timing leak by executing a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths mismatch to strictly balance the execution time.
