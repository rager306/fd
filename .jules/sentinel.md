## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-09-13 - Timing Attack in Authentication Middleware
**Vulnerability:** The `APIKeyAuth` middleware used `subtle.ConstantTimeCompare` on slices of unequal lengths, which immediately returns false and leaks the length of the expected secret.
**Learning:** `subtle.ConstantTimeCompare` is only secure if the two slices being compared have exactly the same length. Comparing arbitrary user input directly against a secret using this function still exposes the secret's length through timing side-channels.
**Prevention:** To safely compare arbitrary-length secrets without leaking lengths, hash both the expected secret and the user input (e.g., using `crypto/sha256.Sum256`) and securely compare the resulting fixed-length hashes. Pre-compute the expected hash if possible.
