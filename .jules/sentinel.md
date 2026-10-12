## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-10-12 - Ignored crypto/rand.Read error
**Vulnerability:** Ignored error from crypto/rand.Read(buf) in queue.NewRequestID() which can result in an all-zero ID if random generation fails.
**Learning:** If the random generation fails (e.g., due to entropy depletion), the buffer remains uninitialized (zeroes). Using this for identifiers results in predictable values, leading to collisions or vulnerabilities.
**Prevention:** Always check and handle the error returned by crypto/rand.Read(buf). Panic if random generation is strictly required for security.
