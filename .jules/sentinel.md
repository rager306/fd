## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-26 - Length-based Timing Leaks in ConstantTimeCompare
**Vulnerability:** Length-based timing leak in API key authentication because `subtle.ConstantTimeCompare` returns early when the two slices have different lengths.
**Learning:** Passing raw strings of variable lengths to `subtle.ConstantTimeCompare` creates a timing side-channel that leaks the exact length of the expected secret.
**Prevention:** Hash the expected secret once during initialization and hash the incoming input on each request before passing them to `ConstantTimeCompare`, ensuring both slices have identical, fixed lengths (e.g. 32 bytes for SHA-256).
