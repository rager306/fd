## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack in API Key Authentication
**Vulnerability:** The APIKeyAuth middleware used `subtle.ConstantTimeCompare` directly on strings (the token and API key) of potentially different lengths. Since `ConstantTimeCompare` returns immediately if the lengths differ, this creates a timing side-channel that leaks the expected length of the API key.
**Learning:** `ConstantTimeCompare` must only be used on fixed-size slices, or slices where the length difference is already known and doesn't leak sensitive information.
**Prevention:** To securely compare secrets of potentially variable lengths (like API tokens), hash both inputs (e.g., using `crypto/sha256.Sum256`) and compare the resulting fixed-length hashes.
