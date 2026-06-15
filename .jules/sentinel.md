## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-24 - Missing Security Headers
**Vulnerability:** The HTTP API was missing standard security headers (X-Content-Type-Options, X-Frame-Options, Strict-Transport-Security).
**Learning:** These headers were omitted during initial global middleware setup.
**Prevention:** Ensure new services or endpoints are bootstrapped with a standard set of security headers globally.
