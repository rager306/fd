## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Length-Based Timing Attack in ConstantTimeCompare
**Vulnerability:** The APIKeyAuth middleware used `subtle.ConstantTimeCompare` on unhashed inputs of potentially different lengths.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ. This allows an attacker to deduce the length of the expected secret (like an API key) via timing attacks.
**Prevention:** Always hash secrets/tokens first (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to guarantee equal-length inputs.
