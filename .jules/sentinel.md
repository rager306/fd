## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-30 - Timing Side-Channel in ConstantTimeCompare
**Vulnerability:** Comparing API tokens directly with `subtle.ConstantTimeCompare` exposed the expected token's length because the function returns immediately if lengths differ.
**Learning:** Go's `subtle.ConstantTimeCompare` only protects against timing attacks for the contents of byte slices, not their lengths. Exposing the length reduces brute-force search space.
**Prevention:** For comparing secrets of variable lengths, hash both the expected secret and the provided token first, then compare the fixed-length hashes.
