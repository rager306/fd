## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-10-18 - Length-Based Timing Attack in ConstantTimeCompare
**Vulnerability:** Length-based timing attack during API key validation in `api/middleware/auth.go`.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ, which can leak the expected length of a secret.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee equal length.
