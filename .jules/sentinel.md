## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-07-19 - Timing Attack in API Key Comparison
**Vulnerability:** The API key middleware used `subtle.ConstantTimeCompare` directly on the raw request token and expected API key, which leaks the key length because `ConstantTimeCompare` returns early if lengths differ.
**Learning:** `subtle.ConstantTimeCompare` only provides constant-time comparison when the inputs are of the same length. Comparing variable-length secrets directly re-introduces timing attacks.
**Prevention:** Always hash secrets (e.g., with `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to guarantee equal-length inputs. Calculate the expected hash outside the HTTP handler closure to avoid redundant computations.
## 2026-07-19 - HTTP/3 QPACK Vulnerability in quic-go
**Vulnerability:** The project was using `github.com/quic-go/quic-go@v0.59.0` which is vulnerable to GO-2026-5676 (HTTP/3 QPACK Trailer Expansion Memory Exhaustion). This allows an attacker to exhaust server memory when handling HTTP/3 requests due to inefficient trailer decompression.
**Learning:** `govulncheck` accurately tracks call graphs to determine if a project actually executes vulnerable paths inside dependency modules rather than just identifying metadata presence. Here, `http.Server.ListenAndServe` and `CacheHandler.Delete` hit the vulnerable HTTP/3 `quic-go` routines.
**Prevention:** Regularly scan backend applications using `govulncheck ./...` in CI pipelines and bump dependencies proactively using `go get module@version` when a CVE fix is released (e.g. `v0.59.1`).
