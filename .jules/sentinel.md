## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-27 - Length-based Timing Attack in API Key Validation
**Vulnerability:** The APIKeyAuth middleware used `subtle.ConstantTimeCompare` directly on user input and the expected API key. If the lengths mismatched, `ConstantTimeCompare` returned immediately, allowing an attacker to determine the length of the valid API key through a timing attack.
**Learning:** Even when using `subtle.ConstantTimeCompare`, if the input length differs from the expected secret's length, Go's implementation leaks this mismatch in O(1) time. In high-performance hot-paths where cryptographic hashing (like SHA-256) is too expensive, we must manually mask the length mismatch to ensure a constant time execution profile.
**Prevention:** When comparing tokens in a hot path, always perform a length-hiding comparison. If lengths mismatch, assign the expected secret to the comparison buffer so `ConstantTimeCompare` always compares two slices of the exact same length, masking the true length.

## 2024-07-27 - quic-go HTTP/3 QPACK Vulnerability (GO-2026-5676)
**Vulnerability:** The codebase had an outdated dependency `github.com/quic-go/quic-go@v0.59.0` which was vulnerable to HTTP/3 QPACK Trailer Expansion Memory Exhaustion.
**Learning:** Continuous monitoring and updating of underlying networking dependencies (especially for emerging protocols like HTTP/3/QUIC) is critical because protocol parsing logic frequently handles unvalidated external input. Memory exhaustion vectors like this can lead to remote Denial of Service (DoS).
**Prevention:** Regularly scan dependencies using `govulncheck` in CI and promptly upgrade modules to their patched versions (e.g., `github.com/quic-go/quic-go@v0.59.1`) upon detection.
