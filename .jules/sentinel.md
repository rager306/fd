## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-06 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** API key comparison in `middleware/auth.go` used `subtle.ConstantTimeCompare` directly on raw string bytes, which returns immediately if lengths differ, allowing an attacker to deduce the expected API key length via timing attacks.
**Learning:** While `subtle.ConstantTimeCompare` protects against content-based timing attacks, its early-return behavior on length mismatch still leaks the secret's length.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee equal lengths and prevent length-based timing attacks.
