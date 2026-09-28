## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - ConstantTimeCompare Length Leaks
**Vulnerability:** The APIKeyAuth middleware directly compared the incoming bearer token with the expected API key using `subtle.ConstantTimeCompare`. This function requires both slices to be of equal length. Since the incoming token is user-provided, comparing unequal slices immediately reveals that the token length is incorrect, and worse, subtle timing differences in short-circuiting can theoretically be exploited.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time if the lengths of the two byte slices are identical. If they differ, it leaks the length.
**Prevention:** Always hash both the expected secret and the user-provided token (e.g., with SHA-256) before passing them to `subtle.ConstantTimeCompare`. This guarantees both slices have a fixed, identical length (32 bytes).
