## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-18 - API Key Length Timing Leak
**Vulnerability:** The APIKeyAuth middleware leaked the exact length of the API key via timing attack.
**Learning:** `crypto/subtle.ConstantTimeCompare` returns immediately if the two provided byte slices have different lengths.
**Prevention:** Always hash secrets (e.g., `sha256.Sum256`) to guarantee equal length before passing them to `ConstantTimeCompare`.
