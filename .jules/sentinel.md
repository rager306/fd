## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-27 - Timing side-channel in subtle.ConstantTimeCompare
**Vulnerability:** Comparing an attacker-controlled token directly against a secret API key using `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))` is vulnerable to a timing side-channel attack if the length of the tokens differ, as `ConstantTimeCompare` returns immediately if the slices are of different lengths.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time if the inputs are guaranteed to be the same length. It does not mask length differences.
**Prevention:** To securely compare secrets of potentially variable lengths, hash both inputs first using a strong hashing algorithm (like SHA-256) to ensure they are the exact same length (e.g., 32 bytes) before passing them to `ConstantTimeCompare`.
