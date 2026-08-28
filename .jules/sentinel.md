## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-28 - Timing Attack in API Key Validation
**Vulnerability:** API key validation used `subtle.ConstantTimeCompare` directly on raw strings of variable lengths, exposing a timing side-channel that leaks the length of the secret API key.
**Learning:** `ConstantTimeCompare` immediately returns 0 if byte slices have different lengths. This is a common pitfall when comparing variable-length secrets like API tokens.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256.Sum256`) before comparison to ensure both slices are of identical, fixed lengths, preventing length leakage.
