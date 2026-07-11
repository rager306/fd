## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-07-11 - Length-Based Timing Attack in Auth Middleware
**Vulnerability:** `subtle.ConstantTimeCompare` used directly on variable-length API keys/tokens.
**Learning:** `ConstantTimeCompare` returns early if lengths differ, allowing an attacker to guess the secret's length via timing attacks.
**Prevention:** Always hash secrets (e.g., with SHA-256) before passing them to `ConstantTimeCompare` to guarantee equal length.
