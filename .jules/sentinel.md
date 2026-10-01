## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-10-01 - Length-based Timing Oracle in API Key Validation
**Vulnerability:** `subtle.ConstantTimeCompare` was returning immediately when the provided token length did not match the expected API key length. This created a timing side-channel that could leak the length of the secret API key.
**Learning:** Go's `subtle.ConstantTimeCompare` only provides constant-time comparison for slices of the *same* length. It returns 0 immediately if lengths differ.
**Prevention:** When comparing secrets against user input without hashing, always balance execution time by running a dummy comparison (e.g. against the secret itself) if lengths mismatch, before rejecting the input.
