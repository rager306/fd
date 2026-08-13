## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-01-01 - API Key Timing Attack via subtle.ConstantTimeCompare
**Vulnerability:** API key verification was vulnerable to a timing attack because `subtle.ConstantTimeCompare` leaks the length difference in O(1) time if slice lengths do not match.
**Learning:** `ConstantTimeCompare` only guarantees constant time execution if the lengths of the slices are identical. Differing lengths return immediately, leaking the secret length.
**Prevention:** Mask length mismatches by verifying if lengths differ, setting a validity flag to 0, matching the strings (e.g. comparing the secret to itself), and applying a bitwise AND to the constant time comparison result with the validity flag.
