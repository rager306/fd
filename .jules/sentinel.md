## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Attack in API Key Auth
**Vulnerability:** Comparing API keys using `subtle.ConstantTimeCompare` without hashing first allowed for a length-based timing attack because the function returns early if lengths differ.
**Learning:** Even when using "constant time" functions, you must ensure the inputs are of equal length, otherwise the early length check defeats the purpose of the constant-time operation.
**Prevention:** When comparing secrets of potentially varying lengths, first hash them (e.g., with `crypto/sha256`) to ensure equal length before passing them to `subtle.ConstantTimeCompare`.
