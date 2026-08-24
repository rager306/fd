## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-25 - Timing side-channel in API Key validation
**Vulnerability:** `subtle.ConstantTimeCompare` leaked the length of the secret API key because it returns early when slices have different lengths.
**Learning:** Even when using constant-time comparison functions, comparing variable-length secrets (like API tokens) directly can expose the expected secret's length.
**Prevention:** Always hash variable-length secrets (e.g., using `crypto/sha256.Sum256`) before comparing them with `subtle.ConstantTimeCompare`, ensuring both inputs are fixed-length arrays. Pre-calculate the expected secret's hash in middleware initialization to avoid overhead.
