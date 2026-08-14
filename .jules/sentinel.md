## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-08-14 - Dependency Update to fix Vulnerability
**Vulnerability:** HTTP/3 QPACK Trailer Expansion Memory Exhaustion in `github.com/quic-go/quic-go`.
**Learning:** `govulncheck` accurately detected standard library and module vulnerabilities. It is crucial to address reported vulnerabilities in CI effectively by updating dependencies rather than just silencing lint checks.
**Prevention:** Regularly scan with `govulncheck` and update to secure versions, specifically `github.com/quic-go/quic-go` to `v0.59.1`.
