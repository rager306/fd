## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-03 - Prevent Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on variable-length API tokens, exposing the service to length-based timing attacks because it returns early if lengths differ.
**Learning:** Even when using timing-safe comparison functions, input lengths must match to guarantee constant-time execution. Otherwise, attackers can iteratively determine the secret length.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to ensure equal lengths.
