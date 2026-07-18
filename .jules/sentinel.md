## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-18 - Length-Based Timing Attack in API Key Auth
**Vulnerability:** API key verification used `subtle.ConstantTimeCompare` directly on string bytes. If lengths differed, it returned early, allowing timing attacks to discover the length of the valid key.
**Learning:** `subtle.ConstantTimeCompare` is only constant time if both byte arrays have the same length. Comparing secrets of variable lengths exposes the length.
**Prevention:** Always hash secrets (e.g., using `sha256.Sum256`) before comparing them with `subtle.ConstantTimeCompare` to ensure the comparison always takes constant time regardless of the original inputs.
