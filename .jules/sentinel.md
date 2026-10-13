## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-13 - Length-Based Timing Oracle in Token Comparison
**Vulnerability:** Comparing an arbitrary length token and secret API key using `subtle.ConstantTimeCompare` immediately returns 0 if lengths differ, leaking the length of the secret key.
**Learning:** `ConstantTimeCompare` acts as a timing oracle when string lengths differ.
**Prevention:** Hash both strings (e.g., using SHA-256) before passing them to `ConstantTimeCompare` to guarantee uniform lengths and prevent length-based timing attacks.
