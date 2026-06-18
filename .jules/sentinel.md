## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-14 - Fix authorization length-based timing attack
**Vulnerability:** Length-based timing attack in API Key validation where `subtle.ConstantTimeCompare` compared plain text strings.
**Learning:** `subtle.ConstantTimeCompare` returns early if the slices have different lengths, which leaks the expected key length to the attacker.
**Prevention:** Hash both strings using `crypto/sha256` before passing them to `subtle.ConstantTimeCompare` to guarantee equal length slices.
