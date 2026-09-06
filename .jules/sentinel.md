## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - Timing Attack in API Key Comparison
**Vulnerability:** The APIKeyAuth middleware used `subtle.ConstantTimeCompare` directly on strings of potentially varying lengths. This exposes a timing side-channel that leaks the length of the expected secret because `ConstantTimeCompare` returns 0 immediately if lengths differ.
**Learning:** Go's `subtle.ConstantTimeCompare` is only safe to use on slices of identical lengths.
**Prevention:** Always hash secrets of potentially variable lengths (like API tokens) to a fixed length (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare`.
