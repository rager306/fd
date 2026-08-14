## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-27 - Constant-Time Comparison Timing Attack
**Vulnerability:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, leaking the length of the secret or the provided token in O(1) time.
**Learning:** In high-performance hot-paths without pre-hashing, direct use of `subtle.ConstantTimeCompare` exposes a timing side-channel.
**Prevention:** Mask length mismatches by reassigning the byte slice to evaluate securely in constant time, utilizing bitwise AND against a length-validity flag, and verify byte length using `len(string)` directly to avoid unnecessary allocations.
