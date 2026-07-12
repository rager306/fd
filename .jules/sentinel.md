## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-03-01 - Prevent Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** Length-based timing attack in API key validation because `subtle.ConstantTimeCompare` returns early if lengths differ.
**Learning:** Even when using `subtle.ConstantTimeCompare`, if the input lengths are not guaranteed to be equal, the function returns immediately, leaking the length of the secret. Hashing the inputs first guarantees equal lengths.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare`.
