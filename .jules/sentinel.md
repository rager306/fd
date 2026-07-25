## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-25 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** Comparing plaintext tokens and API keys using `subtle.ConstantTimeCompare` leaks the expected API key length via early return on mismatch.
**Learning:** Even when using constant-time comparison functions, comparing strings or byte slices of different lengths will leak the length of the expected secret because the function immediately returns false if the lengths differ.
**Prevention:** Hash both the user-provided token and the expected secret using a cryptographic hash function (like SHA-256) before comparison. This ensures both inputs have the same fixed length, preventing length-based timing attacks.
