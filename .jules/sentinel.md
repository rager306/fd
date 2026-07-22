## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-05-27 - Timing Attack in API Key Comparison
**Vulnerability:** The `APIKeyAuth` middleware used `subtle.ConstantTimeCompare` directly on the provided token and expected API key strings.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices differ. An attacker could use timing analysis to deduce the exact length of the expected API key, aiding in brute-force or side-channel attacks.
**Prevention:** When comparing secrets of potentially varying lengths, hash both inputs using a cryptographically secure hash function (e.g., SHA-256) before using `ConstantTimeCompare`. This guarantees both inputs have the same length, ensuring the comparison time is truly constant and preventing length-based timing attacks.
