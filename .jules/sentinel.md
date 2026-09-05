## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-25 - API Key Length Leak
**Vulnerability:** The `APIKeyAuth` middleware used `subtle.ConstantTimeCompare` directly on byte slices of potentially different lengths. This function returns 0 immediately if the lengths differ, exposing a timing side-channel that leaks the secret length.
**Learning:** `ConstantTimeCompare` only provides constant-time properties if the inputs are of the same length. Variable length string comparisons need to be normalized to a fixed length before comparison.
**Prevention:** Always hash secrets to a fixed length (e.g. SHA-256) before using `ConstantTimeCompare`. Compute the expected hash once at initialization to prevent per-request overhead.
