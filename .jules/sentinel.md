## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-02 - Length-Based Timing Attack in API Key Validation
**Vulnerability:** API key validation was vulnerable to length-based timing attacks because `subtle.ConstantTimeCompare` returns immediately if the input byte slices are of different lengths.
**Learning:** Comparing arbitrary length secrets directly using `ConstantTimeCompare` leaks the length of the expected secret.
**Prevention:** Ensure both inputs are first hashed using a cryptographically secure hash function (e.g., `crypto/sha256`) to guarantee equal length before calling `ConstantTimeCompare`.
