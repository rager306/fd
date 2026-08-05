## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - API Key length leakage in middleware
**Vulnerability:** Length of the API key was leaked in an O(1) time timing attack because `subtle.ConstantTimeCompare` returns immediately if slice lengths differ.
**Learning:** `subtle.ConstantTimeCompare` is only constant time when the slices have the same length. It does not mask differences in length, which is a common pitfall.
**Prevention:** If input lengths can differ, manually mask the mismatch by temporarily replacing the differing slice with the expected secret to force constant-time evaluation, then logically `AND` the result with a validity flag.
