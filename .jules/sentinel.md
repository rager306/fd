## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-04 - Authentication Timing Attack via subtle.ConstantTimeCompare
**Vulnerability:** The authentication middleware was using `subtle.ConstantTimeCompare` to verify the bearer token, but `subtle.ConstantTimeCompare` returns immediately in O(1) time if the lengths of the two byte slices are different. This can leak the length of the secret API key.
**Learning:** `subtle.ConstantTimeCompare` only guarantees constant time execution if the lengths are equal.
**Prevention:** In hot paths where hashing both inputs is too expensive, mask the length mismatch. If the lengths differ, compare the expected secret against itself, and use a bitwise AND with a length-validity flag against the comparison result.
