## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Prevent Timing Attack in ConstantTimeCompare
**Vulnerability:** Length-based timing leak in API key validation because `subtle.ConstantTimeCompare` returns immediately if lengths differ.
**Learning:** `subtle.ConstantTimeCompare` only provides constant time execution *if the lengths are equal*. When comparing user input against secrets where the secret length should be protected, this creates a timing oracle that allows an attacker to guess the correct length of the secret.
**Prevention:** When comparing secrets where length could vary, perform a dummy constant-time comparison against the secret itself if the lengths do not match to balance execution time.
