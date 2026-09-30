## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-01 - API Key Length Timing Attack
**Vulnerability:** The APIKeyAuth middleware leaked the exact length of the API key due to `subtle.ConstantTimeCompare` returning early on length mismatch.
**Learning:** Even when using constant-time comparison functions, input lengths must be handled carefully. If the function short-circuits on unequal lengths, it creates a timing oracle that allows an attacker to brute-force the length of the secret.
**Prevention:** When comparing secrets against user input without hashing, execute a dummy constant-time operation when lengths mismatch to balance the execution time.
