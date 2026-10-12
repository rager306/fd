## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-12 - Timing Attack Oracle via ConstantTimeCompare
**Vulnerability:** Length-based timing attack in API token validation.
**Learning:** In Go, `subtle.ConstantTimeCompare` acts as a timing oracle because it immediately returns 0 if slice lengths differ, leaking the length of the secret key to attackers testing different string lengths.
**Prevention:** Always hash sensitive strings of potentially variable length (like API keys) using a secure hashing algorithm (like SHA-256) before passing them to `ConstantTimeCompare` to guarantee fixed uniform lengths.
