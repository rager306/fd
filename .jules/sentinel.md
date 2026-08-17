## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-17 - API Key Timing Attack
**Vulnerability:** API key verification used `subtle.ConstantTimeCompare` directly on variable-length tokens.
**Learning:** `ConstantTimeCompare` immediately returns 0 if byte slices have different lengths, exposing the length of the secret API key through a timing side-channel.
**Prevention:** To securely compare secrets of variable lengths, hash both inputs (e.g., `sha256.Sum256`) to a fixed length before using `ConstantTimeCompare`.
