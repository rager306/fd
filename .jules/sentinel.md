## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - Timing Attack via ConstantTimeCompare
**Vulnerability:** subtle.ConstantTimeCompare returns 0 immediately if the slices are of different lengths, exposing a timing side channel that leaks the length of the expected secret.
**Learning:** When comparing secrets of variable lengths, directly using subtle.ConstantTimeCompare creates a vulnerability because of the length-check short circuit.
**Prevention:** Always compare the fixed-length cryptographic hashes (like SHA-256) of the secrets instead of the secrets themselves when using subtle.ConstantTimeCompare.
