## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-19 - Length-based Timing Attack in API Key Auth
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the provided `token` and the expected `apiKey`.
**Learning:** `subtle.ConstantTimeCompare` leaks the length of the inputs if they are different lengths (it returns early). If the token and API key are directly compared, an attacker can determine the length of the secret API key by sending tokens of varying lengths and measuring the response time.
**Prevention:** When comparing secrets or tokens using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length.
