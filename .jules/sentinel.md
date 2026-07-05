## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - API Key Length-Based Timing Attack
**Vulnerability:** `subtle.ConstantTimeCompare` was comparing raw tokens, allowing an attacker to determine the length of the configured API key via timing side-channels since `ConstantTimeCompare` returns early if lengths differ.
**Learning:** When comparing secrets of variable or unknown length, `ConstantTimeCompare` is insufficient on its own because it leaks length information.
**Prevention:** Always hash the secrets (e.g., using `crypto/sha256`) before passing them to `ConstantTimeCompare` to ensure both inputs have a fixed, identical length.
