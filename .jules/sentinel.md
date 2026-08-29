## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-29 - Timing Attack in API Key Auth
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on byte slices of potentially different lengths. Because this function immediately returns `0` if lengths differ, it exposed a timing side-channel that leaked the length of the secret API key.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time if the inputs are already known to be the exact same length. If lengths differ, it short-circuits.
**Prevention:** To securely compare secrets of variable or unknown lengths (like API tokens), hash both inputs (e.g., using `crypto/sha256.Sum256`) and compare the resulting fixed-length hashes.
