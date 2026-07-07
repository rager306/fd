## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-07 - Length-Based Timing Attack in ConstantTimeCompare
**Vulnerability:** Comparing plaintext secrets using `subtle.ConstantTimeCompare` without hashing first.
**Learning:** `subtle.ConstantTimeCompare` returns early if lengths differ, allowing an attacker to determine the secret's length through timing analysis.
**Prevention:** Both inputs must be hashed (e.g., using `crypto/sha256`) prior to comparison to guarantee equal length execution and fully prevent timing attacks.
