## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-17 - Length-Based Timing Attack via ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` returned immediately if the compared byte slices had different lengths, leaking the expected secret length via timing differences.
**Learning:** `ConstantTimeCompare` is only constant-time for inputs of the *same* length. When comparing variable-length secrets (like API tokens), failing to hash them first exposes the secret's length.
**Prevention:** Securely compare secrets of potentially variable lengths by hashing both inputs (e.g., using `crypto/sha256.Sum256`) and comparing the resulting fixed-length hashes.
