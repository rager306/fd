## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-22 - Timing Attack in API Key Validation
**Vulnerability:** `APIKeyAuth` middleware used `subtle.ConstantTimeCompare` directly on the raw bearer token and the API key. `subtle.ConstantTimeCompare` returns early if the slices have different lengths, which exposes the validation to length-based timing attacks where an attacker could theoretically guess the expected length of the API key.
**Learning:** When using `subtle.ConstantTimeCompare` to compare secrets like tokens or API keys, it's essential to first hash the inputs (e.g., with `crypto/sha256.Sum256`) to guarantee equal length and prevent length-based timing attacks.
**Prevention:** Always ensure both inputs to `subtle.ConstantTimeCompare` have been hashed to the same fixed-size length before comparison.
