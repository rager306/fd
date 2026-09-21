## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack on Auth Middleware
**Vulnerability:** The APIKeyAuth middleware directly compared incoming tokens with the expected apiKey using `subtle.ConstantTimeCompare` with byte slices of potentially different lengths.
**Learning:** Using `subtle.ConstantTimeCompare` on slices of different lengths allows a length-based timing attack, as the comparison fails immediately if the lengths differ.
**Prevention:** Always hash incoming tokens and expected secrets using a robust cryptographic hash (like SHA-256) before comparison to guarantee fixed-length inputs to `subtle.ConstantTimeCompare`.
