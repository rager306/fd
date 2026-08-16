## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-16 - API Key Length Timing Attack
**Vulnerability:** The API key verification used `subtle.ConstantTimeCompare` directly on the provided token and the expected API key. If their lengths differed, it returned 0 immediately, allowing attackers to guess the length of the API key via a timing attack.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time for inputs of the same length. It returns immediately if the lengths differ.
**Prevention:** When comparing secrets where length could be sensitive or unknown (like tokens or passwords), hash both values using a strong cryptographic hash function (like SHA-256) first, and then compare the fixed-length hashes using `ConstantTimeCompare`.
