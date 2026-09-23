## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-24 - Length-Based Timing Attack in API Key Check
**Vulnerability:** API key verification in authentication middleware used `subtle.ConstantTimeCompare` without hashing.
**Learning:** `subtle.ConstantTimeCompare` terminates early if the byte slices have different lengths. Attackers can guess the API key length by observing timing differences.
**Prevention:** Always hash secrets (e.g., with SHA-256) before using `ConstantTimeCompare` to ensure the byte slices always have fixed, identical lengths, preventing length leakage. Calculate the expected secret's hash once at initialization for performance.
