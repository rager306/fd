## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-16 - Prevent length-based timing attacks in subtle.ConstantTimeCompare
**Vulnerability:** Length-based timing attack in `subtle.ConstantTimeCompare` inside `api/middleware/auth.go`.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the inputs differ, which exposes the exact length of the API key through timing side-channels.
**Prevention:** Hash both the expected key and the provided token using `crypto/sha256` prior to comparison, guaranteeing uniform lengths. Additionally, computing the expected hash outside the middleware closure optimizes performance.
