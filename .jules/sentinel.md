## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-01-01 - ConstantTimeCompare Timing Attack
**Vulnerability:** Length-based timing attack in `api/middleware/auth.go` where `subtle.ConstantTimeCompare` leaks the length of the secret API key by returning early if the token length is different.
**Learning:** Go's `subtle.ConstantTimeCompare` requires inputs of the same length to run in constant time. A simple length check mitigates this, but to completely avoid leaking timing data about the expected secret length in hot-paths where pre-hashing is expensive, one should compare the secret against itself on a length mismatch.
**Prevention:** Ensure length masking (comparing the secret against itself) is used in combination with `subtle.ConstantTimeCompare` in hot-paths, or use pre-hashing for generic secret comparisons where performance allows.
