## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-20 - Prevent Length-Based Timing Attacks in API Key Validation
**Vulnerability:** `subtle.ConstantTimeCompare` was used to compare raw API keys of potentially different lengths, which allows early exit length-based timing attacks.
**Learning:** `ConstantTimeCompare` returns immediately if the lengths of the two inputs differ. This leaks the expected length and fails to protect against timing attacks when length is variable. Precomputing the hash outside the closure also improves throughput.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to guarantee they are the same length.
