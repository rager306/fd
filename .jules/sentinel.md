## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - API Key Length Leak
**Vulnerability:** The APIKeyAuth middleware compared the incoming bearer token with the expected API key using `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`.
**Learning:** `subtle.ConstantTimeCompare` returns early if the two slices have different lengths. This allowed an attacker to enumerate the length of the secret API key by sending tokens of varying lengths and observing the timing difference.
**Prevention:** Always hash the expected secret and the incoming token (e.g., with SHA-256) before passing them to `subtle.ConstantTimeCompare`. This guarantees both slices have the same fixed length (e.g., 32 bytes), preventing early returns and length leaks.
