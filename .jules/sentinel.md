## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-14 - Timing Attack in API Key Comparison
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the token and expected API key strings. `ConstantTimeCompare` returns early if lengths differ, exposing a timing side-channel that leaks the length of the secret key.
**Learning:** For variable-length secrets (like API keys), direct comparison with `ConstantTimeCompare` is insufficient because the length check is not constant-time. Both inputs must be hashed to a fixed length (e.g., using `crypto/sha256`) before comparison to ensure completely constant-time execution regardless of input length. Also, we can pre-compute the expected key's hash once during middleware initialization to avoid overhead on every request.
**Prevention:** Always hash variable-length secrets before comparing them with `ConstantTimeCompare`.
