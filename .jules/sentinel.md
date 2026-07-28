## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-28 - HTTP/3 QPACK Trailer Expansion Memory Exhaustion
**Vulnerability:** `github.com/quic-go/quic-go` up to `v0.59.0` is vulnerable to GO-2026-5676, allowing attackers to exhaust server memory via crafted QPACK trailers in HTTP/3 connections.
**Learning:** Found by checking CI annotations resulting from `govulncheck` scans identifying the dependency flaw. The issue is severe (memory exhaustion DoS) on any internet-facing endpoint routing over HTTP/3.
**Prevention:** Keep network listener dependencies like `quic-go` updated proactively by automating `govulncheck` pipelines and bumping module versions as soon as patches (like `v0.59.1`) are published.
