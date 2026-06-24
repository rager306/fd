## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-14 - Fix length-based timing attack in API key validation
**Vulnerability:** API key verification using `subtle.ConstantTimeCompare` with directly casted `[]byte(token)` and `[]byte(apiKey)` allowed length-based timing attacks.
**Learning:** `subtle.ConstantTimeCompare` performs a constant-time comparison *only* if the lengths of the two inputs are equal. If the lengths differ, it returns immediately in a non-constant time manner, allowing an attacker to deduce the length of the secret API key.
**Prevention:** Hash both the input token and the expected secret using a cryptographic hash function (like `crypto/sha256`) *before* applying `ConstantTimeCompare`. This guarantees both inputs have the identical length.
