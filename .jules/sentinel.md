## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-01-01 - HTTP/3 QPACK Trailer Expansion Memory Exhaustion
**Vulnerability:** The application used an outdated version of `github.com/quic-go/quic-go` (v0.59.0) which had a known memory exhaustion vulnerability (CVE/GO-2026-5676) due to QPACK trailer expansion.
**Learning:** Outdated network dependencies can introduce critical vulnerabilities like DoS through memory exhaustion, even if the application code itself is secure.
**Prevention:** Regularly scan dependencies with tools like `govulncheck` and update them promptly to patch known CVEs.
