## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-10-14 - Length-Based Timing Attack in API Key Validation
**Vulnerability:** The API key validation used `subtle.ConstantTimeCompare` directly on the plaintext API key and incoming token.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two inputs differ. This leaks the exact length of the expected API key via timing analysis, weakening the secret.
**Prevention:** Always hash secrets (e.g., with `crypto/sha256`) before passing them to `ConstantTimeCompare` to guarantee inputs are of equal length, neutralizing length-based timing attacks. In HTTP middleware, the expected secret should be hashed once outside the handler closure for performance.
