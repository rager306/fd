## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-27 - [Length-based timing attack in subtle.ConstantTimeCompare]
**Vulnerability:** Length-based timing attack due to unequal length inputs in `subtle.ConstantTimeCompare`.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ. This can allow an attacker to determine the length of the API key by sending requests with tokens of varying lengths and observing the response times.
**Prevention:** Always hash the inputs to guarantee equal length before comparing them with `subtle.ConstantTimeCompare`.
