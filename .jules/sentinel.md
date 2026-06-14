## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-14 - Timing Attack via crypto/subtle
**Vulnerability:** Timing attack allowing length-leakage of API keys in auth middleware.
**Learning:** In Go, `subtle.ConstantTimeCompare` only provides constant-time comparison if the two byte slices have the exact same length. If lengths differ, it returns 0 immediately, leaking the length of the secret.
**Prevention:** When comparing secrets of potentially varying lengths against user input, always hash both strings using a strong cryptographic hash (like SHA-256) first. The hashes will always be of fixed length, ensuring `ConstantTimeCompare` evaluates in constant time.
