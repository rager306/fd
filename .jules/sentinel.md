## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - ConstantTimeCompare Length Leak
**Vulnerability:** subtle.ConstantTimeCompare returns immediately when slice lengths differ, leaking the length in O(1) time.
**Learning:** Go's subtle.ConstantTimeCompare is only constant-time if the lengths match, requiring manual masking for variable-length secrets.
**Prevention:** Check if lengths differ, set a flag to 0, reassign variables to mask the mismatch, and bitwise AND the result to ensure constant-time evaluation.
