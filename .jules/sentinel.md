## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-03-01 - Timing Attack in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` returned early on length mismatch, leaking the length of the API key via timing side-channels.
**Learning:** Go's `subtle.ConstantTimeCompare` requires both byte slices to be of equal length. Otherwise, it exits immediately, negating the constant-time guarantee and leaking length information.
**Prevention:** To securely compare secrets of variable length, compute a fixed-length hash (e.g., SHA-256) of both the expected secret and the provided input, then use `ConstantTimeCompare` on the resulting hashes. Compute the hash of static secrets at initialization to avoid per-request overhead.
