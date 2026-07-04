## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** API key comparison used `subtle.ConstantTimeCompare` directly on plaintext tokens, exposing the length of the secret.
**Learning:** `subtle.ConstantTimeCompare` returns early if input lengths differ, allowing an attacker to determine the secret's length via timing attacks.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) to guarantee equal length before comparing them with `ConstantTimeCompare`.
