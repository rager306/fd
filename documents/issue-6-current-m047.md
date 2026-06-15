# Audit follow-up — G4/G5: dependency resilience, graceful shutdown, error contract

Follow-up to #3. The M046 remediation closed the P0/P1 disclosure-risk cluster (default-open auth, batch-endpoint DoS, N+1 TEI, LocalCache races). This issue tracks the **still-open G4 (dependency-failure resilience) + G5 (error/shutdown contract)** findings plus `getEnvInt` overflow — the remaining P2/P3 items worth a focused follow-up. Descriptions, not fixes.

> All locations re-verified against current `master` (HEAD `a16e243`). Build/vet/tests are green; these are reliability/contract gaps, not compile or test failures.

---

## G4 — Dependency-failure resilience

### #11 — TEI HTTP client has no retry / backoff / circuit breaker  · P2
**Where:** `api/embed/tei.go:79` — `resp, err := c.httpClient.Do(req); if err != nil { return nil, err }`.
**What breaks:** TEI is the service's only critical dependency. Any transient TEI failure (container restart, 503 during model swap, network blip) turns into an immediate `500 internal_error` with no recovery. The 30s HTTP timeout (main.go) bounds hanging but does not retry. For a service whose entire job is calling TEI, this is the core reliability surface.
**Now:** confirmed — `retry|backoff|circuit` absent from `tei.go`.
**Suggested fix:** wrap the call in a bounded retry loop (2–3 attempts) with jittered exponential backoff for retriable failures (network errors, 502/503/504); add a lightweight circuit breaker that short-circuits to `503 overloaded_error` after N consecutive failures, so the handler does not burn the full 30s timeout on every call during a TEI outage.

### #14 — Warmup failure permanently degrades readiness; no auto-retry  · P2
**Where:** `api/main.go:129-135` — on `PreWarm` error, `state.SetLastError(err)` is called and the goroutine returns; no retry.
**What breaks:** `State.IsReady()` returns false while `LastError() != nil`, and `LastError` is only cleared by `MarkWarmupDone`. A single warmup failure at boot (TEI still loading, 5s `defaultWarmupTimeout`) leaves `/v1/embeddings` and `/ready` returning `503 model_not_loaded` indefinitely — until a manual `POST /warmup` or a process restart. For a dependency that may itself still be starting (container orchestration ordering), one 5s timeout is a fragile readiness gate.
**Now:** confirmed — no auto-retry/backoff for warmup.
**Suggested fix:** on warmup failure, schedule a retried warmup with exponential backoff (e.g. 2s→4s, capped) up to a bounded number of attempts before giving up; `MarkWarmupDone` clears the error on success.

---

## G5 — Error & shutdown contract

### #13 — `ListenAndServe` fatal error calls `os.Exit(1)`, bypassing graceful shutdown  · P2
**Where:** `api/main.go:295-297` — `if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { ... os.Exit(1) }`.
**What breaks:** a fatal server error (port conflict surfacing after start, accept-loop exhaustion) calls `os.Exit(1)` immediately — skipping `BeginShutdown()` and the 30s `GracefulShutdown` drain. In-flight embedding requests are not tracked as draining, the Redis pool is not closed (`closeResource` is unreachable), and callers get connection resets instead of `503 shutting_down`. The signal path handles shutdown correctly; the server-error path does not.
**Now:** confirmed — `main.go:295,297` unchanged.
**Suggested fix:** instead of `os.Exit(1)`, signal the main goroutine (error channel or `cancel()` on a context main selects on) so the same `GracefulShutdown` drain path runs; minimum — call `lifecycleState.BeginShutdown()` and a bounded drain before exiting.

### #32 — `errorlint`: direct sentinel comparison instead of `errors.Is`  · P3
**Where:** `api/main.go:295` — `err != http.ErrServerClosed` (same line as #13).
**What breaks:** violates `.golangci.yml` `errorlint` (binding per M043); a wrapped `ErrServerClosed` would mis-classify as a fatal error. Pre-existing.
**Suggested fix:** `!errors.Is(err, http.ErrServerClosed)` + import `errors`. (Natural to fold into the #13 change since both touch this line.)

### #25 — Error codes registered but never emitted  · P3
**Where:** `api/handlers/errors.go:17,28,29` — `CodeDimensionsRequired`, `CodeDimensionsMismatch`, `CodeRequestTimeout` are declared, in `errorCodeRegistry`, and in `AllErrorCodes()` (which seeds Prometheus zero-series), but no handler emits them.
**What breaks:** premature API surface — codes defined for behaviors never implemented (dimensions-mismatch detection, request timeouts). They occupy error-catalog rows and Prometheus cardinality but can never be produced; a maintainer may wire logic against codes that have no emitter.
**Now:** confirmed — still present.
**Suggested fix:** remove the three codes (and their registry/`AllErrorCodes` rows), or annotate `reserved-for-future-use` with the specific behavior that will emit them, and add a test asserting every registered code is emitted by at least one non-test path.

---

## Bonus (bundled — same "resilience/contract" follow-up)

### #15 — `getEnvInt` overflow silently disables the in-flight capacity gate  · P2
**Where:** `api/main.go:52-65` — hand-rolled `n = n*10 + int(c-'0')` loop with no overflow guard; feeds `maxInFlight := getEnvInt("FD_MAX_IN_FLIGHT", 0)` (main.go:216).
**What breaks:** a malformed/overlong numeric `FD_MAX_IN_FLIGHT` overflows int and can silently disable the `LifecycleGateWithCapacity` gate (or behave unpredictably). No error surfaced. (The `envInt` copies in `ratelimit.go` and `lru.go` already use `strconv.Atoi`; only `main.go`'s `getEnvInt` is still hand-rolled.)
**Suggested fix:** replace the loop with `strconv.Atoi` (errors on overflow → keep default) and reject `n < 0`.

---

## Scope note

Not in this issue: the P3 cleanup tail (`#19` dead `LRUCache`, `#24` `handleBindError` message, `#26` ONNX dead health fields, `#27` dup hash helpers, `#28` `envInt` dedup, `#30` `defaultState` singleton, `#31` `openapi.m()` silent drop) and `#7` (`/metrics` public — likely accepted for monitoring), `#22` (OpenAPI omits `/embeddings/batch`), `#16` (`/v1/traces` scoping — mitigated by fail-closed auth, single-tenant by design). These can be a separate housekeeping pass.

This issue documents the problems only — implementation to be proposed/PR'd separately.
