## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-28 - Timing Side-Channel in API Key Auth
**Vulnerability:** The API key middleware used `subtle.ConstantTimeCompare` directly on the provided token and the expected API key. `subtle.ConstantTimeCompare` returns immediately if lengths differ, allowing attackers to guess the secret's length.
**Learning:** Even when using "constant-time" comparison functions, variable-length inputs can introduce timing leaks if the function short-circuits on length mismatches.
**Prevention:** Always hash variable-length secrets (like API keys) using a secure hash function (e.g., SHA-256) before comparing them with `subtle.ConstantTimeCompare` to ensure inputs are fixed-length.
