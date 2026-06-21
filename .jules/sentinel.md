## 2026-06-21 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** API key verification used `subtle.ConstantTimeCompare` directly on raw string bytes of unequal lengths, allowing a length-based timing attack since the function returns early if slice lengths differ.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time properties if the input slices are of equal length. Otherwise, the comparison takes variable time proportional to the length check.
**Prevention:** When comparing secrets of potentially different lengths, hash both inputs first (e.g., `sha256.Sum256`) to guarantee equal length before calling `ConstantTimeCompare`.

## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
