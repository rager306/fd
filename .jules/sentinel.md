## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-23 - Timing Attack in API Key Validation
**Vulnerability:** The API key validation middleware used `subtle.ConstantTimeCompare` directly on the incoming token and the expected API key. If the strings have different lengths, this can leak the length of the expected key or allow timing attacks.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time if the slices being compared are of the same length. Comparing variable-length strings directly exposes length information.
**Prevention:** Always hash both the expected secret and the user-provided token (e.g., using SHA-256) before passing them to `subtle.ConstantTimeCompare` to guarantee fixed-length inputs.
