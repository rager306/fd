## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was being called directly on the raw incoming bearer token and the configured API key byte slices. If the lengths differed, it would return early, enabling a potential length-based timing attack.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time comparison when the lengths of both inputs are identical. When dealing with user-supplied input strings, lengths can vary.
**Prevention:** When comparing secrets or tokens using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length before comparison.
