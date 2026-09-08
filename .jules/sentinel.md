## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - Timing Attack in API Key Validation
**Vulnerability:** Timing side-channel attack in API key validation via `subtle.ConstantTimeCompare` which leaks the secret length.
**Learning:** `subtle.ConstantTimeCompare` immediately returns 0 if the byte slices have different lengths, exposing a timing side-channel. To securely compare secrets of potentially variable lengths, we must hash both inputs and compare the resulting fixed-length hashes.
**Prevention:** Always hash secrets before comparing them with `subtle.ConstantTimeCompare` to ensure length is always identical. Hash the static secret during initialization to avoid per-request overhead.
