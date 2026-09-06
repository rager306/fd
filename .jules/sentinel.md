## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-25 - Timing Attack in Secret Comparison
**Vulnerability:** The API authentication middleware compared a variable-length token against a secret using `subtle.ConstantTimeCompare`. Because this function checks length first and returns early, it leaks the secret's length via a timing side-channel.
**Learning:** Even when using cryptographic comparison functions, passing variable-length, attacker-controlled input directly can expose side-channel vulnerabilities.
**Prevention:** For comparing secrets of potentially variable lengths (like API tokens), always hash both the expected secret and the incoming input (e.g., via SHA-256) to fixed lengths before comparing them. Hashing static secrets should happen once during initialization.
