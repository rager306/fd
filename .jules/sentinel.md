## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-24 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** The API key verification used `subtle.ConstantTimeCompare` directly on raw string bytes. This function returns early if lengths differ, allowing attackers to guess the API key length via timing attacks.
**Learning:** `ConstantTimeCompare` is only constant-time for inputs of equal length. Passing arbitrary-length strings directly exposes length information.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before using `ConstantTimeCompare` to ensure both slices are of equal, fixed length.
