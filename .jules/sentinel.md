## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-25 - Length-Based Timing Attack in Token Comparison
**Vulnerability:** Comparing user-provided tokens with API keys using `subtle.ConstantTimeCompare` without hashing first.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ. This exposes the length of the expected secret via a timing attack, which can make brute-forcing easier or bypass some protections.
**Prevention:** Always ensure both inputs are equal in length before using `ConstantTimeCompare`. The most robust way to do this is to hash both the secret and the user input (e.g., using `crypto/sha256`) and compare the resulting hashes.