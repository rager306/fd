## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-17 - Length-based Timing Attack in Auth Middleware
**Vulnerability:** The APIKeyAuth middleware was susceptible to length-based timing attacks due to short-circuiting the comparison of unequal length inputs before hitting `subtle.ConstantTimeCompare`.
**Learning:** `subtle.ConstantTimeCompare` leaks information if the two slices are of unequal length because it immediately returns 0 in that case, bypassing constant-time comparison entirely. This reveals the exact length of the expected secret to an attacker.
**Prevention:** To safely use `subtle.ConstantTimeCompare` with untrusted input lengths, hash both the expected secret and the provided input (e.g., using `crypto/sha256`) and compare their fixed-length hashes instead.
