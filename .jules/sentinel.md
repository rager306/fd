## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Attack in Auth Middleware
**Vulnerability:** The APIKeyAuth middleware directly compared the incoming Bearer token and the expected API key using `subtle.ConstantTimeCompare`. However, if the lengths do not match, this function returns 0 immediately, leaking the length of the expected API key to an attacker via a timing difference.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time comparison when the lengths of the two byte slices are identical. Passing variable-length inputs can lead to length-based timing attacks.
**Prevention:** Hash both the expected secret and the incoming token (e.g., using SHA-256) before comparison. This ensures fixed-length slices. The static expected secret should be hashed once during middleware initialization to avoid re-computing overhead on every request.
