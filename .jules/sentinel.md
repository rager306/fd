## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Predictable Request ID Generation
**Vulnerability:** The `newRequestID` function fell back to using the current time (`time.Now().UnixNano()`) to generate a predictable Request ID when `crypto/rand` failed.
**Learning:** Request IDs are used as correlation keys, and in the queue subsystem (`queue_handlers.go`), they are used as identifiers to retrieve embedding results. If an attacker can predict a queued task's ID, they could exploit an Insecure Direct Object Reference (IDOR) to access another user's processing data.
**Prevention:** Never use predictable pseudo-randomness (like current time) when system entropy (`crypto/rand`) fails, especially when identifiers have security implications. Always "fail securely" by panicking or returning an error.
