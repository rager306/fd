## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-11 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on the bearer token and API key, allowing an attacker to determine the API key length via a timing attack.
**Learning:** Comparing plaintext secrets directly with `ConstantTimeCompare` leaks the length of the expected secret because it returns early if lengths differ.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) to guarantee equal lengths before comparing them with `subtle.ConstantTimeCompare`.
