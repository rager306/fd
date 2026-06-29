## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-01-20 - Prevent length-based timing attacks in ConstantTimeCompare
**Vulnerability:** Length-based timing attack via early exit in `subtle.ConstantTimeCompare` for unequal length inputs.
**Learning:** `ConstantTimeCompare` returns immediately if slice lengths differ, which can leak the length of the expected API key or secret.
**Prevention:** Always hash secrets (e.g., with SHA-256) prior to `ConstantTimeCompare` to guarantee identical slice lengths.
