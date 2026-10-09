## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-15 - Addressing CI govulncheck Failures
**Vulnerability:** Known CVEs in HTTP/2 (`golang.org/x/net`) and HTTP/3 (`github.com/quic-go/quic-go`) caused CI failures via `govulncheck`. Specifically, memory exhaustion and flow control bugs (GO-2026-6612, GO-2026-6611, GO-2026-6603, and GO-2026-5676).
**Learning:** `govulncheck` runs in CI as an explicit standalone check (`golang.org/x/vuln/cmd/govulncheck@latest`) and will hard-fail the build if dependencies have known CVEs that are reachable in the codebase.
**Prevention:** Always run `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` locally if CI includes it, and resolve flagged issues by bumping dependencies (e.g., `go get golang.org/x/net@latest`) before submitting a PR.
