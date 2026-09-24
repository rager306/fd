## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack via ConstantTimeCompare Lengths
**Vulnerability:** The authentication middleware compared incoming tokens directly with the expected API key using `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`.
**Learning:** `ConstantTimeCompare` leaks the length of the slices because it immediately returns if the lengths differ. This allows an attacker to guess the length of the secret API key through timing side-channels, weakening the overall security.
**Prevention:** Always hash secrets (e.g., with SHA-256) before passing them to `ConstantTimeCompare`. This guarantees fixed-length inputs, eliminating length-based timing leaks.
