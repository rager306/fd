## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-10-05 - Length-based timing attack in API key verification
**Vulnerability:** API key verification was vulnerable to a length-based timing attack because `subtle.ConstantTimeCompare` returns immediately if lengths differ.
**Learning:** Returning immediately allows an attacker to brute force the length of the secret API key.
**Prevention:** If the provided string and the correct API key length mismatches, perform a dummy comparison `subtle.ConstantTimeCompare(apiKey, apiKey)` to balance execution time.
