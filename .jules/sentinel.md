## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-18 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on plaintext tokens of arbitrary length, which could allow a timing attack if tokens are of unequal length, as `ConstantTimeCompare` checks length first and returns early.
**Learning:** `subtle.ConstantTimeCompare` does not compare strings of different lengths in constant time. It only works in constant time for inputs of the same length.
**Prevention:** When comparing a user-provided secret (like a token or API key) with a stored secret, always hash both inputs (e.g., with `crypto/sha256`) before passing them to `ConstantTimeCompare`. This guarantees both inputs will be exactly the same length.
