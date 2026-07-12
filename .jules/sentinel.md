## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-03-01 - Dependency Vulnerability Mitigation (GO-2026-5856, GO-2026-5676)
**Vulnerability:** The application was vulnerable to GO-2026-5856 (an Encrypted Client Hello privacy leak in `crypto/tls` from the Go standard library) and GO-2026-5676 (an HTTP/3 QPACK Trailer Expansion Memory Exhaustion in `github.com/quic-go/quic-go`).
**Learning:** Standard library and deeply nested dependencies can introduce security flaws over time. Dependency vulnerability scanners (like `govulncheck`) are vital for proactively identifying these before they can be exploited.
**Prevention:** Regularly scan dependencies using `govulncheck` in CI and promptly update vulnerable packages to their patched versions.
