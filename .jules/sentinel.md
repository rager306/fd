## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-08-06 - ConstantTimeCompare Length Leakage
**Vulnerability:** Timing attack via early return in `subtle.ConstantTimeCompare` when comparing `token` against `apiKey`.
**Learning:** `subtle.ConstantTimeCompare` leaks the length of inputs because it checks lengths first and returns early in O(1) time if they differ, exposing the length of the secret key.
**Prevention:** In high-performance hot-paths, mask the length difference. Track the mismatch with a valid flag (1 or 0), copy the secret bytes over the input bytes to make lengths equal, and perform a bitwise AND on the result of `subtle.ConstantTimeCompare` with the valid flag to ensure constant-time evaluation.
