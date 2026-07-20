## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-10-24 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ, which can lead to a timing attack disclosing the expected length.
**Learning:** Comparing tokens directly with `subtle.ConstantTimeCompare` without hashing can be vulnerable to length-based timing attacks.
**Prevention:** Always hash both the expected token and the provided token (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to ensure equal lengths.
