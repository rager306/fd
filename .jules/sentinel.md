## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-10-11 - Length-Based Timing Oracle in Token Comparison
**Vulnerability:** API key comparison in `auth.go` using `subtle.ConstantTimeCompare` without length checks leaked the secret length via timing differences.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, creating a timing oracle that attackers can use to deduce the secret's length.
**Prevention:** Always perform a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths differ before rejecting the token, ensuring strictly balanced execution time.
