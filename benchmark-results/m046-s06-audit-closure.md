# M046 S06 Audit Closure Matrix

Captured: 2026-06-14

## Scope

Close M046 audit remediation for GitHub issue #3 after S01-S05 by:

- fixing residual P1 #6 for `/v1/embeddings` cache peeks;
- fixing P1 #9 405 error contract drift;
- recording a finding-by-finding closure matrix for all 32 findings;
- identifying residual deferred work for future milestones.

Issue snapshot: `documents/issue-3-current-m046.md`.

## S06 Code Changes

- `/v1/embeddings` now calls `EmbeddingCache.GetManyIfPresent(ctx, chunk, dims)` once per bounded chunk.
- `TieredCache.GetManyIfPresent` checks L1 for all inputs, performs one Redis `MGET` for L2 misses, backfills L1 on Redis hits, and returns hits by original input index.
- `RedisCache.GetMany` decodes MGET binary payloads and omits misses, malformed values, and dimension mismatches.
- `CodeMethodNotAllowed` is now registered and `MethodNotAllowedMiddleware` uses `WriteError` instead of a hardcoded envelope.

## Verification

### Red Evidence

```bash
cd api && go test ./handlers -run TestCreateEmbeddingUsesBatchedCachePeek
```

Initial result:

```text
GetManyIfPresent calls = 0, want 1
```

### Final Gates

```bash
cd api && go test ./...
cd api && go test -race ./cache -run TestLocalCache
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
```

Results:

```text
go test ./...: 284 passed in 9 packages
go test -race ./cache -run TestLocalCache: 9 passed in 1 package
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
```

Static proof:

- `gsd_exec:43c16c32-c290-499a-a42a-b8602a0ce6ee` — `/v1/embeddings` uses batched cache peek and Redis MGET.
- `gsd_exec:ab9ac45b-a646-4b6e-a5ef-22839e715e5c` — final proof for batched cache peek and canonical 405.

## Closure Matrix

| # | Priority | Finding | Status | Evidence / Rationale |
|---|---|---|---|---|
| 1 | P0 | Empty `FD_API_KEY` disabled auth | Fixed | S04: protected endpoints fail closed without key; `/v1/embeddings` runtime UAT returns 401. Evidence `benchmark-results/m046-s04-exposure-posture.md`. |
| 2 | P0 | `/embeddings/batch` mounted bare | Fixed | S02: route now has body cap, rate limit, lifecycle/capacity gate. Evidence `benchmark-results/m046-s02-batch-guardrails.md`. |
| 3 | P0 | `/v1/batch` skipped validation/rate/input caps | Fixed | S02: route guardrails plus handler-specific count/length validation. Evidence `benchmark-results/m046-s02-batch-guardrails.md`. |
| 4 | P1 | `/embeddings/batch` N+1 TEI calls | Fixed | S03: legacy batch uses bounded miss chunking. Evidence `benchmark-results/m046-s03-batch-backend-chunking.md`. |
| 5 | P1 | `/v1/batch` N+1 TEI calls | Fixed | S03: v1 batch uses one embedder call per bounded inner batch/miss group. Evidence `benchmark-results/m046-s03-batch-backend-chunking.md`. |
| 6 | P1 | `/v1/embeddings` sequential Redis GET cache peek | Fixed | S06: `GetManyIfPresent` + Redis MGET per chunk. Evidence `43c16c32-c290-499a-a42a-b8602a0ce6ee` and final gates. |
| 7 | P1 | `/metrics` auth-exempt | Fixed | S04: `/metrics` removed from public carve-outs; runtime UAT returns 401 without key. |
| 8 | P1 | Rate limiter spoofable/unbounded map | Fixed | S04: `SetTrustedProxies(nil)` by default plus `maxRateLimitKeys` pruning/eviction. Evidence `benchmark-results/m046-s04-exposure-posture.md`. |
| 9 | P1 | 405 `method_not_allowed` absent from registry/bypassed `WriteError` | Fixed | S06: `CodeMethodNotAllowed` registered and `MethodNotAllowedMiddleware` uses `WriteError`. Final tests and static proof pass. |
| 10 | P1 | `LocalCache` size-counter race/eviction drift | Fixed | S05: mutex-owned map, derived size, idempotent `Close`, race tests. Evidence `benchmark-results/m046-s05-localcache-correctness.md`. |
| 11 | P2 | TEI HTTP client lacks retry/backoff/circuit breaker | Deferred | Reliability hardening remains useful but is outside P0/P1 exposure/DoS/cache correctness wave. Future milestone should design retries around idempotency, timeout budgets, and overload semantics. |
| 12 | P2 | TieredCache swallows Redis errors; L2 loss invisible | Deferred | Not fixed in M046. S06 added MGET warning on batched get failure, but broader Redis error metrics/health belong with observability/resilience follow-up. |
| 13 | P2 | `ListenAndServe` `os.Exit(1)` bypasses graceful drain | Deferred | Not fixed in M046. Requires lifecycle/server error-channel redesign beyond audit P0/P1 wave. |
| 14 | P2 | Warmup failure permanently degrades readiness; no auto-recovery | Deferred | Not fixed in M046. Should be handled with warmup retry/backoff and explicit readiness semantics in a resilience milestone. |
| 15 | P2 | `getEnvInt` overflow can disable in-flight gate | Deferred | Not fixed in M046. Existing tests cover normal parsing; overflow hardening is small and should be included in residual cleanup. |
| 16 | P2 | `/v1/traces` exposes request metadata | Fixed by S04 auth posture | `/v1/traces` is not public after S04 fail-closed auth; only probes/docs/OpenAPI remain public. |
| 17 | P2 | Redis `protected-mode no` / bind posture | Accepted with constraint | Compose does not publish Redis in base; override binds Redis to `127.0.0.1:6379`. Same-host single-tenant contract remains. Future deployment hardening can revisit Redis auth/protected-mode. |
| 18 | P2 | LocalCache O(N) scan on every Set at capacity | Mitigated / accepted | S05 removed `sync.Map` scan/counter drift. Capacity enforcement still evicts by iterating a map when over cap, acceptable for current 10k L1 and outside P1 correctness issue. True LRU can be future work. |
| 19 | P2 | `LRUCache` dead production code | Accepted / future cleanup | LRU remains useful for tests/alternate cache surface and now implements `GetManyIfPresent`. Not removed in M046 to avoid nonessential churn. |
| 20 | P2 | Validation duplicated between middleware and handler | Accepted with rationale | S02 intentionally kept handler-specific validation for batch JSON shapes while sharing generic body cap. Full unification would risk over-abstraction. |
| 21 | P2 | Endpoints validate inconsistently | Mostly fixed | S02 added shared batch limits and body cap with shape-specific validation. Full schema abstraction deferred as nonessential after P0/P1 guardrails pass. |
| 22 | P2 | OpenAPI omits `/embeddings/batch` | Deferred | Legacy endpoint remains undocumented in OpenAPI; adding contract docs for legacy compatibility is useful but outside M046 P0/P1 remediation. |
| 23 | P3 | LocalCache eviction goroutine/ticker leak | Fixed | S05 added idempotent `Close()` and main shutdown integration. |
| 24 | P3 | `handleBindError` malformed array-element error | Deferred | Not fixed in M046; should be included in API error polish follow-up. |
| 25 | P3 | Dead error catalog entries | Accepted | Error registry intentionally includes canonical codes used across lifecycle and docs; S06 added `method_not_allowed` rather than pruning catalog. |
| 26 | P3 | `RuntimeHealth` ONNX-only dead fields | Accepted with context | ONNX is inactive/future research; health model remains safe TEI runtime metadata. Avoided churn in M046. |
| 27 | P3 | Duplicate hash-truncate helpers | Deferred | Maintainability cleanup only; no security/correctness impact in M046. |
| 28 | P3 | `envInt` duplicated / hand-rolled ASCII loop | Deferred | Low-risk cleanup; future refactor can centralize env parsing. |
| 29 | P3 | `Embedder` and `WarmupModel` same interface under two names | Accepted | Names document package-level roles and avoid cross-package coupling. No M046 action. |
| 30 | P3 | `lifecycle.defaultState` singleton redundant | Deferred | Low-risk maintainability cleanup. |
| 31 | P3 | `openapi.m()` silently drops non-string keys | Deferred | Low-risk OpenAPI helper cleanup; future docs/spec pass. |
| 32 | P3 | `err != http.ErrServerClosed` should use `errors.Is` | Deferred | Small standards cleanup; not included in M046 after final gates passed. |

## Requirement Outcomes

- R029 validated by S02/S03: batch endpoints bounded and backend work shaped.
- R030 validated by S04: protected endpoints fail closed and exposure posture is safer.
- R031 validated by S05: LocalCache accounting/lifecycle correctness.
- R032 validated by S06: `/v1/embeddings` cache peeks use batched lookup.

## Residual Follow-up Candidates

Recommended future slices/milestones, in priority order:

1. Dependency resilience: TEI retry/backoff/circuit breaker, Redis error metrics, warmup retry (#11/#12/#14).
2. API contract polish: 405/OpenAPI/validation array error and legacy `/embeddings/batch` spec docs (#22/#24 plus regression docs).
3. Maintainability cleanup: env parsing, duplicate helpers, lifecycle singleton, optional LRU cleanup (#19/#27/#28/#30/#31/#32).

## Verdict

M046 closes the audit's P0 and P1 findings (#1-#10). Remaining P2/P3 findings are either fixed as side effects, accepted under current same-host constraints, or deferred with explicit follow-up direction.
