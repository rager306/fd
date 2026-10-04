## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-12 - Timing Oracle in subtle.ConstantTimeCompare
**Vulnerability:** The API Key authentication middleware had a timing attack oracle. `subtle.ConstantTimeCompare` returns immediately if slice lengths differ, allowing an attacker to infer the expected API key length based on response times.
**Learning:** Even when using cryptographic constant-time functions, length checks are still susceptible to timing leaks. In modern Go, string to byte slice conversions passed directly to `ConstantTimeCompare` do not allocate, but hashing (which would eliminate the length difference) can introduce DoS vulnerabilities.
**Prevention:** Always ensure balanced execution time for sensitive comparisons by executing a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths mismatch, rather than returning early or relying on hashing.
