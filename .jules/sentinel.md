## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-28 - HTTP/3 QPACK Trailer Expansion Memory Exhaustion
**Vulnerability:** GO-2026-5676 in `github.com/quic-go/quic-go` version `v0.59.0`
**Learning:** `github.com/quic-go/quic-go` < `v0.59.1` is vulnerable to memory exhaustion via QPACK trailer expansion.
**Prevention:** Monitor dependencies regularly using `govulncheck` and update to `v0.59.1` or later to fix HTTP/3 QPACK memory exhaustion vulnerabilities.
