## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Timing Attack in API Key Comparison
**Vulnerability:** Used `subtle.ConstantTimeCompare` directly on raw string inputs of potentially variable lengths, leaking the secret length via early return.
**Learning:** Go's `subtle.ConstantTimeCompare` immediately returns 0 if byte slices have different lengths, exposing a timing side-channel.
**Prevention:** Always hash secrets (e.g., using `crypto/sha256.Sum256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee fixed-length inputs.
