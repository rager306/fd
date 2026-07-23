## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Length-based timing attack in Auth
**Vulnerability:** Length-based timing attack in token comparison due to direct usage of `subtle.ConstantTimeCompare` with dynamic length inputs.
**Learning:** `subtle.ConstantTimeCompare` only protects against timing attacks if both byte slices have the same length; it returns immediately if the lengths differ, which leaks the expected length.
**Prevention:** Hash both inputs (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to ensure both have a fixed, identical length.
