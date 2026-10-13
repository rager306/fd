## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-10-13 - Unhandled crypto/rand Error
**Vulnerability:** The queue ID generation logic using `crypto/rand.Read` ignored the returned error.
**Learning:** If `crypto/rand.Read` fails (e.g., due to entropy depletion), it leaves the buffer uninitialized. This results in predictable identifiers, which could lead to critical vulnerabilities like IDOR or authentication bypasses.
**Prevention:** Always verify and handle the error returned by `crypto/rand.Read`, panicking if random generation is strictly required for security and cannot fallback.
