## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-10-11 - Length-based Timing Oracle in API Key Validation
**Vulnerability:** API key validation used `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`. `ConstantTimeCompare` returns immediately if the lengths differ, leaking the length of the expected API key via timing side-channels.
**Learning:** `ConstantTimeCompare` is only fully constant-time if the inputs are guaranteed to be the same length. Comparing arbitrary-length secrets directly exposes their lengths.
**Prevention:** Always hash sensitive strings (e.g., using SHA-256) before passing them to `ConstantTimeCompare` so the lengths are always identical and fixed.
