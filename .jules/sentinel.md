## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-06 - Prevent Length-Based Timing Attacks in API Key Auth
**Vulnerability:** Length-based timing attack oracle in API Key validation.
**Learning:** `crypto/subtle.ConstantTimeCompare` returns immediately if input slice lengths do not match, causing execution time to be proportional to string length, allowing attackers to deduce expected lengths.
**Prevention:** Always compare against the secret's own length when inputs differ in length using a dummy comparison `subtle.ConstantTimeCompare(secret, secret)` to ensure consistent execution time.
