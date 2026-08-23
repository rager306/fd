## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - API Key Timing Attack
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the provided token and the expected API key. This function returns immediately if the lengths differ, exposing a timing side-channel that leaks the length of the expected secret.
**Learning:** To securely compare secrets of potentially variable lengths, they must first be hashed to a fixed length before comparison.
**Prevention:** Hash both inputs (e.g., using `crypto/sha256.Sum256`) and compare the resulting fixed-length hashes using `subtle.ConstantTimeCompare`. Calculate the hash of static expected secrets (like environment API keys) once during middleware initialization to avoid unnecessary overhead.
