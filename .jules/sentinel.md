## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - HTTP/3 QPACK Trailer Expansion Memory Exhaustion
**Vulnerability:** A vulnerability in `github.com/quic-go/quic-go` (GO-2026-5676) could cause memory exhaustion and a potential DoS attack by sending crafted HTTP/3 requests, leading to server crashes.
**Learning:** We rely on `quic-go` for HTTP/3, and these vulnerabilities propagate. Even though our code might not seem directly related, vulnerabilities deep in the network stack affect us.
**Prevention:** Keep dependencies updated, actively monitor vulnerability scanners like `govulncheck`, and run them in CI to catch vulnerabilities early.
