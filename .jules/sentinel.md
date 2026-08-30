## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Side-Channel in API Key Validation
**Vulnerability:** `subtle.ConstantTimeCompare` leaks the length of secrets when the byte slices being compared have different lengths, as it immediately returns `0` on length mismatch. This exposes a timing side-channel.
**Learning:** When comparing variable-length secrets (like tokens or API keys provided in headers) against an expected secret, you cannot compare them directly with `subtle.ConstantTimeCompare`.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256.Sum256`) to ensure fixed-length byte slices before comparing them with `subtle.ConstantTimeCompare`. For performance, pre-compute the hash of static expected secrets.
