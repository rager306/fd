## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-10 - Prevent length-based timing attacks in token comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on tokens, returning early if lengths differ, allowing length-based timing attacks.
**Learning:** `ConstantTimeCompare` does not hide the length of inputs, which can leak the expected API key length.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee equal lengths.
