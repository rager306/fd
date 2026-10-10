## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-10 - Length-based timing oracle in token comparison
**Vulnerability:** The API key validation leaks the correct API key length through early exit upon length mismatch. `subtle.ConstantTimeCompare` returns immediately if slice lengths differ.
**Learning:** When comparing secrets against user input without hashing, length-based timing oracles must be mitigated by balancing execution time even when lengths mismatch.
**Prevention:** Use a dummy comparison `subtle.ConstantTimeCompare(secret, secret)` when lengths mismatch to strictly balance execution time.
