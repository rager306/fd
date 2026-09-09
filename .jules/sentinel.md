## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-28 - Timing Attack in API Key Validation
**Vulnerability:** Go's `subtle.ConstantTimeCompare` returns immediately if the byte slices have different lengths. Using it to directly compare a user-provided token against a secret API key creates a timing side-channel that leaks the length of the secret.
**Learning:** When comparing variable-length secrets (like API keys), hashing both the expected secret and the user input ensures both sides are a fixed length, eliminating length-based timing leaks. Pre-computing the hash of static secrets in middleware initialization avoids per-request overhead.
**Prevention:** Always use hashing (like SHA-256) when comparing variable-length sensitive data to prevent timing side-channels, even when using constant-time comparison functions.
