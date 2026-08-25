## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Side-Channel in API Key Comparison
**Vulnerability:** The APIKeyAuth middleware used `subtle.ConstantTimeCompare` directly on byte slices of potentially differing lengths.
**Learning:** `ConstantTimeCompare` returns immediately if lengths do not match, creating a timing side-channel that leaks the length of the secret API key.
**Prevention:** Always hash both inputs before comparing them using `ConstantTimeCompare`, ensuring fixed-length inputs and preventing length leakage.
