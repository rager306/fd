## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-01-22 - Length-Based Timing Attack in API Key Auth
**Vulnerability:** API key verification in `api/middleware/auth.go` passed plain byte slices of differing lengths to `subtle.ConstantTimeCompare`, returning early if lengths didn't match and allowing a timing attack to guess token lengths.
**Learning:** `ConstantTimeCompare` only provides constant-time evaluation if the two byte slices are the exact same length. If lengths differ, it fails fast, leaking the length information.
**Prevention:** Always hash secret tokens (e.g. using `crypto/sha256`) before passing them to `ConstantTimeCompare`. This guarantees both inputs have a fixed, identical length, regardless of the underlying token sizes.
