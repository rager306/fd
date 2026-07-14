## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-18 - Fix Length-Based Timing Attack in API Key Comparison
**Vulnerability:** Length-based timing attack vulnerability in `subtle.ConstantTimeCompare` due to unhashed, variable-length inputs.
**Learning:** `subtle.ConstantTimeCompare` can leak the length of the expected secret because it returns early if the lengths of the two byte slices differ.
**Prevention:** Always hash secrets (e.g., with `crypto/sha256`) to ensure fixed-length inputs before passing them to `ConstantTimeCompare`.
