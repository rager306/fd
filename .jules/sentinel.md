## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-20 - ConstantTimeCompare Length-Based Timing Attack
**Vulnerability:** Length-based timing attack in API Key validation due to using `subtle.ConstantTimeCompare` on strings of different lengths.
**Learning:** `subtle.ConstantTimeCompare` only protects against timing attacks when the two slices are the *same length*. If they differ in length, it returns early, allowing an attacker to determine the expected length of the API key, which weakens its security.
**Prevention:** When comparing secrets or tokens using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length before comparing them.
