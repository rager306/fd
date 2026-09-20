## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-20 - ConstantTimeCompare Length Timing Leak
**Vulnerability:** Comparing API keys of differing lengths using `subtle.ConstantTimeCompare` leaks the length of the expected key.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time if both inputs are identical in length. If they are different lengths, an attacker might deduce the length of the secret token.
**Prevention:** Hash both the expected secret and the user-provided token (e.g., using SHA-256) prior to using `subtle.ConstantTimeCompare`. This guarantees both inputs are a fixed length.
