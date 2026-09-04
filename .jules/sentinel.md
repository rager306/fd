## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing side-channel in subtle.ConstantTimeCompare
**Vulnerability:** subtle.ConstantTimeCompare immediately returns 0 if byte slices have different lengths, exposing a timing side-channel that leaks the length of the expected secret (e.g., API key).
**Learning:** The length of a secret can be leaked through timing variations, which could aid attackers in brute-forcing or deducing the secret structure.
**Prevention:** When comparing secrets of potentially variable lengths, hash both the expected secret (once at init) and the incoming token (per request) to a fixed length (e.g., using SHA-256) before using ConstantTimeCompare.
