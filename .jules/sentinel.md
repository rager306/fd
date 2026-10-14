## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Length-based Timing Oracle in Token Comparison
**Vulnerability:** The API key and bearer token were directly compared using `subtle.ConstantTimeCompare`, which acts as a timing oracle that leaks the length of the secret being compared by returning 0 immediately if the lengths differ.
**Learning:** In Go, `subtle.ConstantTimeCompare` only guarantees constant time execution if the lengths of the two slices are identical. When comparing sensitive strings (like API keys or tokens) where lengths may vary or are unknown, an attacker could potentially infer the length of the configured API key by sending tokens of varying lengths and observing the response times.
**Prevention:** Always hash both strings (e.g., using SHA-256) before passing them to `ConstantTimeCompare` to guarantee uniform lengths and prevent length-based timing attacks.
