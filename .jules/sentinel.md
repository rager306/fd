## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Constant-Time Comparison Timing Attack
**Vulnerability:** `subtle.ConstantTimeCompare` leaked string lengths in O(1) time because it returned early on length mismatch.
**Learning:** Even using a `subtle` package does not automatically provide full timing-attack immunity. The `ConstantTimeCompare` function is only constant-time for slices of identical length. If lengths differ, it short-circuits.
**Prevention:** Always mask length mismatches before comparing by substituting inputs (e.g., compare secret against itself) and applying a bitwise mask to the final result.
