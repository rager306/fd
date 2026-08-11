## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-08-11 - Timing Attack in API Key Length Check
**Vulnerability:** API key length is leaked in O(1) time due to `subtle.ConstantTimeCompare` returning immediately when lengths differ.
**Learning:** `subtle.ConstantTimeCompare` requires inputs of equal length to execute in constant time. Passing unequal lengths breaks the constant-time guarantee.
**Prevention:** Mask length mismatches by comparing the expected secret against itself when lengths differ, and bitwise AND the result with a validity flag.
