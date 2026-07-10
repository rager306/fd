## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-10 - Length-Based Timing Attack
**Vulnerability:** Comparing API keys of varying lengths using `subtle.ConstantTimeCompare` without hashing first.
**Learning:** `subtle.ConstantTimeCompare` can return early if the lengths of the slices differ, leaking length information or enabling timing attacks based on length.
**Prevention:** Always hash secrets (e.g., with `crypto/sha256`) before passing them to `ConstantTimeCompare` to guarantee they have an equal length.
