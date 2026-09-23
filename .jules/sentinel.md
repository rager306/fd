## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-23 - Length-based Timing Vulnerability in Authentication Middleware
**Vulnerability:** The global API key authentication middleware used `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`. Go's `ConstantTimeCompare` returns immediately if the lengths of the slices are different, introducing a length-based timing leak that could allow an attacker to guess the correct length of the expected API key.
**Learning:** Even when using cryptographic comparison functions designed to resist timing attacks, developers must ensure the inputs are of equal length. Otherwise, the function may short-circuit, defeating its purpose.
**Prevention:** When comparing variable-length secrets (like API keys or passwords) using constant-time functions, always hash both the expected secret and the user-provided input (e.g., using SHA-256) before comparison. This guarantees that both inputs have a fixed, identical length, eliminating the timing leak.
