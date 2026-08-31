## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - API Key Length Timing Leak
**Vulnerability:** `subtle.ConstantTimeCompare` leaks the length of the expected secret because it immediately returns 0 if the byte slices have different lengths.
**Learning:** For variable-length secrets (like API tokens), this allows an attacker to brute-force the length of the secret.
**Prevention:** Always hash both the expected secret and the user-provided secret to a fixed length (e.g., using SHA-256) before passing them to `ConstantTimeCompare`.
