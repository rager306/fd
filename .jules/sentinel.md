## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-14 - Constant-Time Length Masking
**Vulnerability:** subtle.ConstantTimeCompare leaks length differences in O(1) time because it returns immediately if lengths do not match.
**Learning:** In high-performance hot-paths where pre-hashing inputs is too expensive, standard constant-time comparison is vulnerable to length-based timing attacks.
**Prevention:** Mask length mismatches by reassigning the user input to match the expected secret's length and tracking validity with a bitwise flag before comparison.
