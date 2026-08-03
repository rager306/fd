## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-18 - Prevent O(1) length leak in subtle.ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, leaking the length of the expected secret in O(1) time. This existed in `api/middleware/auth.go` during API key validation.
**Learning:** While `ConstantTimeCompare` is secure for matching lengths, pre-hashing both inputs is often recommended but can be too expensive for high-frequency hot paths like API middleware.
**Prevention:** To mitigate this timing attack without the cost of cryptographic hashing on every request, check for a length mismatch. If lengths differ, mask the comparison by comparing the expected secret against itself, then bitwise AND the result with the validity flag to ensure the comparison always runs in constant time relative to the expected secret length.
