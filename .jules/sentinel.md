## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-10-17 - Weak Request ID Generation Fallback
**Vulnerability:** `newRequestID` in `api/middleware/headers.go` fell back to a predictable `time.Now().UnixNano()` if `crypto/rand.Read` failed.
**Learning:** Using an uninitialized or predictable buffer for identifiers when entropy is depleted can lead to critical vulnerabilities like IDOR or authentication bypasses.
**Prevention:** Never ignore the error returned by `crypto/rand.Read(buf)`. If random generation is strictly required for security, panic instead of falling back to a weak value.
