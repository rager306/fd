## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-27 - [Length-Based Timing Attack in API Key Comparison]
**Vulnerability:** The API Key authentication middleware (`APIKeyAuth`) was directly passing raw API key slices of potentially varying lengths to `subtle.ConstantTimeCompare`. This exposes a length-based timing attack, as `ConstantTimeCompare` exits early without comparison if the two byte slices differ in length.
**Learning:** `subtle.ConstantTimeCompare` is only constant time *if the slice lengths are identical*. If the lengths are different, it leaks information by immediately returning, allowing attackers to incrementally guess the required key length by measuring response times.
**Prevention:** Always hash secrets/tokens with a cryptographic hash function like `crypto/sha256` *before* passing them to `ConstantTimeCompare`. This normalizes all inputs to the exact same length (e.g., 32 bytes for SHA-256), thereby completely eliminating the length leak.
