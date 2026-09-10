## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-10 - Timing Attack in API Key Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` leaks the length of the expected API key by immediately returning if lengths differ.
**Learning:** Even when using constant-time comparison functions, comparing variable-length strings still leaks their length.
**Prevention:** Always hash secrets of potentially variable lengths (like API tokens or passwords) to a fixed length (e.g. SHA-256) before passing them to `ConstantTimeCompare`.
