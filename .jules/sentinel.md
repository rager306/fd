## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-07-26 - Timing Attack via Length Comparison
**Vulnerability:** In `APIKeyAuth` middleware, `subtle.ConstantTimeCompare` was used to compare the provided token and the configured API key directly. `ConstantTimeCompare` immediately returns 0 if the two slices have different lengths, which exposes the length of the expected API key via a timing side-channel.
**Learning:** Even when using "constant-time" comparison functions, length differences can introduce early returns. This leaks information about the secret, narrowing the search space for brute-force attacks.
**Prevention:** When comparing secrets of potentially different lengths, hash both the expected secret and the provided input using a cryptographic hash function (like SHA-256) *before* passing them to `ConstantTimeCompare`. This guarantees both inputs are exactly the same length (e.g., 32 bytes).
