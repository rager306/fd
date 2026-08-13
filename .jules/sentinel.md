## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - API Key Authentication Timing Attack
**Vulnerability:** The API key validation in `APIKeyAuth` used `subtle.ConstantTimeCompare` without masking length differences, which leaks the API key length in O(1) time and opens up a timing attack vector.
**Learning:** `subtle.ConstantTimeCompare` returns immediately when slice lengths differ. In hot paths, this length leak can be exploited.
**Prevention:** Mask length mismatches before calling `subtle.ConstantTimeCompare`. If lengths differ, set a valid flag to 0 and compare the expected secret against itself, then bitwise AND the flag with the result.
