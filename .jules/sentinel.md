## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-06-27 - Length-Based Timing Attack in Token Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` was used to directly compare two tokens of potentially different lengths. While it compares slices in constant time based on slice length, it returns immediately if the lengths differ, which allows a side-channel timing attack to discern valid API key lengths and structure.
**Learning:** `subtle.ConstantTimeCompare` is only secure against timing attacks when the two inputs are mathematically guaranteed to have the exact same length. If length is variable based on input, an attacker can extract information.
**Prevention:** Always first hash both strings/tokens (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` to ensure both byte slices passed are exactly the same length.
