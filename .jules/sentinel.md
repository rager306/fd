## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-23 - API Key Length-Based Timing Attack
**Vulnerability:** Length-based timing attack in API Key comparison where `subtle.ConstantTimeCompare` was used with slices of potentially different lengths.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the slices differ, making it vulnerable to length-based timing attacks.
**Prevention:** When comparing secrets or tokens using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length.
