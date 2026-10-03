## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-10-03 - [Fix length-based timing oracle in subtle.ConstantTimeCompare]
**Vulnerability:** Length-based timing attack in authentication middleware. `subtle.ConstantTimeCompare` leaks the length of the expected secret because it returns immediately if the lengths of the two slices differ.
**Learning:** Even constant-time comparison functions can leak information if they are not used carefully. In Go, `subtle.ConstantTimeCompare` returns early on length mismatch, which can be exploited to guess the length of the expected secret (API key).
**Prevention:** Always ensure the lengths are checked first and mitigated. To strictly balance execution time when comparing secrets, perform a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths mismatch to hide the early return timing.
