## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-08-18 - API Key Verification Timing Attack
**Vulnerability:** The API key validation in `api/middleware/auth.go` directly compared the incoming Bearer token with the expected API key using `subtle.ConstantTimeCompare`. However, if the byte slices have different lengths, `ConstantTimeCompare` returns immediately without executing the constant-time operation. This exposes a timing side-channel that leaks the length of the expected API key.
**Learning:** `subtle.ConstantTimeCompare` is only constant-time for inputs of the same length. It is not designed to compare inputs where the length of the secret is unknown to the attacker.
**Prevention:** Before using `ConstantTimeCompare` on variables of potentially different lengths, hash both the expected secret and the provided value (e.g., using `crypto/sha256.Sum256`). This ensures both inputs are exactly the same length before comparison, preventing length-based timing attacks.
