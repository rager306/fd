## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-05 - Timing Attack in subtle.ConstantTimeCompare
**Vulnerability:** length-based timing attack in APIKeyAuth middleware. Go's subtle.ConstantTimeCompare leaks length because it returns immediately if slice lengths differ.
**Learning:** Checking subtle.ConstantTimeCompare without checking lengths first makes the function exit early, leaking the expected string length in O(1) time.
**Prevention:** If lengths differ, mask the mismatch by comparing the expected key against itself, and set a boolean mask so the check still fails, but takes constant time.
