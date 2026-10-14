## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-14 - Length-Based Timing Oracle in Token Comparison
**Vulnerability:** API key comparison used `subtle.ConstantTimeCompare` directly on byte slices of varying lengths, acting as a timing oracle.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, leaking the length of the secret API key.
**Prevention:** Always hash sensitive strings (e.g., with SHA-256) before passing them to `ConstantTimeCompare` to guarantee uniform lengths and prevent timing attacks.
