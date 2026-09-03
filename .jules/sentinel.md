## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-09-03 - Timing attack in auth middleware
**Vulnerability:** Go's `subtle.ConstantTimeCompare` returns immediately if the lengths of the two inputs do not match, exposing a timing side-channel that allows an attacker to deduce the length of the configured API key.
**Learning:** `ConstantTimeCompare` only provides constant time properties when inputs are of the same length. Variable-length comparisons are fundamentally vulnerable to length leaks.
**Prevention:** To securely compare secrets of variable lengths, compute fixed-length cryptographic hashes (e.g. SHA-256) of both inputs and compare the hashes using `subtle.ConstantTimeCompare`.
