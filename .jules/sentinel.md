## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-14 - Timing Attack via ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` immediately returns 0 if the byte slices have different lengths. Using it to directly compare a user-provided token with a secret token exposes a timing side-channel that leaks the length of the secret.
**Learning:** Comparing secrets of potentially variable lengths (like API tokens) directly is unsafe even with constant-time functions.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256.Sum256`) and compare the resulting fixed-length hashes to securely compare secrets.
