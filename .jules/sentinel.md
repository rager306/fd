## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-07-25 - Prevent length-based timing attack in subtle.ConstantTimeCompare
**Vulnerability:** subtle.ConstantTimeCompare was comparing plaintext tokens of potentially variable lengths. If lengths do not match, the function returns immediately, leaking length information of the expected secret.
**Learning:** subtle.ConstantTimeCompare is only constant-time for inputs of the exact same length. Comparing arbitrary-length strings like API keys makes the service vulnerable to length-deduction side-channel attacks.
**Prevention:** Always hash secrets (e.g., with SHA-256) before passing them to subtle.ConstantTimeCompare. This guarantees both inputs have the identical, fixed length (e.g., 32 bytes) needed for true constant-time comparison.
