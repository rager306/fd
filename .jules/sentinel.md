## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-07-27 - HTTP/3 QPACK Trailer Expansion Memory Exhaustion
**Vulnerability:** The codebase imported `github.com/quic-go/quic-go` at a vulnerable version (`v0.59.0`), which had a known HTTP/3 QPACK Trailer Expansion Memory Exhaustion vulnerability (GO-2026-5676).
**Learning:** `govulncheck` identified that `http.Server.ListenAndServe` eventually calls a vulnerable code path `http3.ConfigureTLSConfig`. Even if the application isn't actively leveraging HTTP/3 features heavily, loading the compromised module introduces a vector.
**Prevention:** Always ensure core protocol dependencies like `quic-go` are kept up to date and regularly verify them using `govulncheck`. Updated to `v0.59.1` to resolve the issue.
