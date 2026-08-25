## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-27 - Auth Timing Attack
**Vulnerability:** The API key validation logic used `subtle.ConstantTimeCompare` directly on strings. This function exits early if string lengths differ, leaking the length of the expected secret API key through a timing side-channel.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time properties when comparing slices of the *same* length. Using it directly on variable-length inputs (like raw passwords or API keys) re-introduces timing attacks.
**Prevention:** To securely compare variable-length secrets, hash both inputs using a cryptographic hash function (like SHA-256) first, and then run `subtle.ConstantTimeCompare` on the resulting fixed-length hashes.
