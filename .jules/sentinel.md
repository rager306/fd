## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack on API Key Comparison
**Vulnerability:** The API key middleware compared the raw token and expected key using `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`. If lengths differed, it leaked timing information.
**Learning:** `ConstantTimeCompare` is only constant-time if both inputs have the same length. Comparing raw strings of arbitrary length defeats the purpose and introduces length-based timing leaks.
**Prevention:** Hash both the expected secret and the incoming token (e.g., using SHA-256) before comparison. This ensures fixed-length slices. Precompute the expected hash during initialization to save overhead.
