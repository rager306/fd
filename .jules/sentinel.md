## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-25 - Timing Attack in Secret Comparison
**Vulnerability:** Comparing an API token against a secret using `subtle.ConstantTimeCompare` without hashing first allowed a timing attack to leak the secret's length.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices differ. This exposes a side channel for secrets of variable length.
**Prevention:** Hash both inputs (e.g., using `crypto/sha256.Sum256`) to fixed-length hashes before using `subtle.ConstantTimeCompare`, and pre-calculate the expected secret's hash during initialization for performance.
