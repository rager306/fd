## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-27 - Timing Side-Channels in Auth Middleware
**Vulnerability:** The APIKeyAuth middleware compared the expected API key and the provided bearer token directly using `subtle.ConstantTimeCompare`. Because strings/byte slices in Go can have different lengths, `ConstantTimeCompare` returns immediately if the lengths differ, exposing a timing side-channel that leaks the length of the secret API key.
**Learning:** Constant-time comparison functions only protect against character-by-character timing attacks *if* both inputs are guaranteed to be the exact same length. Comparing arbitrary user input (like a token) directly against a secret (like an API key) fails this requirement if their lengths can differ.
**Prevention:** Always hash both the secret and the user input using a secure hashing algorithm (like SHA-256) before passing them to a constant-time comparison function. This guarantees both inputs are a fixed length (e.g., 32 bytes) regardless of the original string lengths.
