## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-19 - Timing Attack in APIKeyAuth
**Vulnerability:** API key comparison used `subtle.ConstantTimeCompare` on plain strings of differing lengths, causing early return and leaking token length.
**Learning:** `ConstantTimeCompare` is only constant time for inputs of the same length. Comparing a variable-length user input with a static secret leaks information about the secret's length and can be exploited.
**Prevention:** Hash both the user input and the expected secret with a cryptographically secure hash function (e.g., SHA-256) before comparison. This ensures both values are fixed length and prevents length-based timing leaks. Calculate the static expected secret's hash outside the request handler.
