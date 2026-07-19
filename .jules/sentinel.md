## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-07-21 - API Key Length Timing Attack
**Vulnerability:** The API key verification middleware used `subtle.ConstantTimeCompare` directly on the provided token and expected API key. If the token and key have different lengths, `ConstantTimeCompare` returns early in Go, allowing an attacker to determine the API key's length via a timing attack.
**Learning:** `subtle.ConstantTimeCompare` in Go does not protect against length-based timing attacks; it only protects against character-by-character timing attacks when lengths match.
**Prevention:** Always hash secrets/tokens to a fixed length (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee equal input lengths and true constant-time execution.
