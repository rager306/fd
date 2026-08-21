## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-25 - Timing Attack via ConstantTimeCompare
**Vulnerability:** Comparing potentially variable-length API keys with `subtle.ConstantTimeCompare` leaks the key length because it returns immediately if lengths differ.
**Learning:** Go's `subtle.ConstantTimeCompare` only protects against timing attacks for the contents of the slices if their lengths are identical. If lengths differ, it short-circuits.
**Prevention:** Always hash both inputs (e.g., `sha256.Sum256`) to create fixed-length slices before using `subtle.ConstantTimeCompare` for secret comparison.
