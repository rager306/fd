## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Oracle in API Key Validation
**Vulnerability:** API key validation used `subtle.ConstantTimeCompare`, which immediately returns if the token lengths mismatch, creating a length-based timing oracle.
**Learning:** Even when using cryptographic constant-time comparison functions, input lengths must be strictly identical or masked to avoid leaking the secret's length through execution timing.
**Prevention:** Execute a dummy comparison (comparing the secret against itself) when lengths mismatch to balance the execution path and mask the early exit.
