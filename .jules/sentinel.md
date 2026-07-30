## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - API Key Timing Attack
**Vulnerability:** The API Key authentication middleware used `subtle.ConstantTimeCompare`, which returns immediately if slice lengths differ, leaking the length of the expected key in O(1) time.
**Learning:** `subtle.ConstantTimeCompare` does not mask length differences. In hot paths where hashing both inputs is too expensive, the length mismatch must be explicitly masked by always executing the comparison on equal-length slices (e.g. comparing the key against itself if lengths differ) and checking the length condition afterward.
**Prevention:** Always verify slice lengths before calling `subtle.ConstantTimeCompare`. If length masking is needed for performance, explicitly implement a constant-time fallback that performs the comparison regardless of the initial length check.
