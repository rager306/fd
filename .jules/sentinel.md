## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-22 - Timing Side-Channel in API Key Comparison
**Vulnerability:** Go's subtle.ConstantTimeCompare returns immediately if the lengths of the two byte slices are different, exposing a timing side-channel that leaks the length of the secret API key.
**Learning:** Even when using "constant-time" comparison functions, length differences can introduce timing leaks if the inputs are not fixed-length.
**Prevention:** Always hash secrets of potentially variable lengths (like API tokens) before comparing them with subtle.ConstantTimeCompare to ensure both inputs have identical lengths.
