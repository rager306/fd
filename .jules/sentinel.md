## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-15 - Timing Attack via Length Short-Circuiting in Go
**Vulnerability:** Length-based timing oracle in authentication middleware caused by `subtle.ConstantTimeCompare` returning immediately on length mismatch.
**Learning:** `subtle.ConstantTimeCompare` is strictly constant-time *only* for slices of equal length. Passing tokens of arbitrary length directly allows attackers to measure validation time to guess valid token lengths.
**Prevention:** Always verify lengths match before invoking `subtle.ConstantTimeCompare`. To avoid leaking information via a conditional short-circuit branch, perform a dummy comparison (e.g., `subtle.ConstantTimeCompare(secret, secret)`) when lengths mismatch to strictly balance execution time.
