## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack in API Token Comparison
**Vulnerability:** API key verification used `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`.
**Learning:** `subtle.ConstantTimeCompare` returns early if the two byte slices have different lengths, which exposes a timing side-channel that allows an attacker to deduce the exact length of the expected secret (e.g., API key) by measuring response times.
**Prevention:** To securely compare variable-length secrets, hash both the expected key and the provided token (e.g., with `crypto/sha256.Sum256`), then use `subtle.ConstantTimeCompare` on the resulting fixed-length hashes (32 bytes each).
