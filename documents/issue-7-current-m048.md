# Audit cleanup tail (P3): dead code, duplicate helpers, type-safety gaps

Follow-up to #3 and #6. The P0/P1 cluster and all of G4/G5 are closed and verified (build green, tests green). This issue tracks the remaining **P3 cleanup/housekeeping tail** — eight low-severity maintainability findings, all re-verified as still present on current `master` (HEAD `433b01d`). No exploits, no bugs — pure tech debt. Descriptions, not fixes.

---

## #19 — `LRUCache` is dead production code  · P3
**Where:** `api/cache/lru.go:32` (`LRUCache`, `NewLRUCache`, `NewLRUCacheFromEnv`, `EmbeddingCacheKey`, `LRUCacheMetrics`).
**What:** A complete second in-memory LRU cache (own TTL, eviction, metrics, env-config) — ~182 lines + tests — coexisting with `LocalCache`. `main.go` wires **only** `LocalCache` (via `TieredCache`); the only non-test reference to `LRUCache` is `fd_v2_cache_integration_test.go` as scaffolding. The env vars `FD_CACHE_SIZE` / `FD_CACHE_TTL_HOURS` that `NewLRUCacheFromEnv` reads are silently inert. Two competing "L1 cache" mental models; a maintainer touching cache behavior must determine which is live.
**Now:** confirmed present.
**Fix:** delete `lru.go` (+ `lru_test.go`, `lru_rapid_test.go`); in `fd_v2_cache_integration_test.go` swap `NewLRUCache` for a `LocalCache`-backed `EmbeddingCache` stub. Restore an LRU-with-metrics only if a future backend needs it (YAGNI today).

## #24 — `handleBindError` emits a garbled message for array-element type errors  · P3
**Where:** `api/middleware/validation.go:141` — `fmt.Sprintf("input[%s] must be string, got %s", typeErr.Field, typeErr.Value)`.
**What:** When the client sends a non-string array element (e.g. `input:[123]`), `json.UnmarshalTypeError.Field` is empty, so the rendered message reads `input[ must be string, got 123` — malformed (stray `[`, no index). The 400 still fires, but the error body is ugly/misleading for clients debugging a bad payload.
**Now:** confirmed present (unchanged).
**Fix:** guard the empty-Field case — emit `input must be an array of strings, got <Value>` when `typeErr.Field == ""`, keep the indexed form otherwise; or validate element types explicitly in the custom `UnmarshalJSON`.

## #26 — `RuntimeHealth` carries 8 ONNX-only fields that are never populated  · P3
**Where:** `api/handlers/health.go:26-35` — `ArtifactID`, `MaxSequenceLength`, `ValidatedMaxSequenceLength`, `ArtifactVerified`, `TokenizerVerified`, `RuntimeLibraryVerified`, `Provider` (commented "set for ONNX only").
**What:** After the TEI-only refactor (D047–D049), the sole construction site (`main.go`) populates only `Backend`/`Model`/`Dimensions`/`ProductionDefault`/`CacheNamespace`; none of the ONNX fields are ever set. The Go type advertises ONNX as a live backend when it isn't — `/health` lies about what the service reports.
**Now:** confirmed present.
**Fix:** remove the seven ONNX-only fields (or move them into a separate `OnnxRuntimeHealth` restored only if ONNX returns). Keep `Backend`, `Model`, `Dimensions`, `ProductionDefault`, `CacheNamespace`.

## #27 — Duplicate hash-truncate helpers in the same package  · P3
**Where:** `api/cache/tiered.go:114` `shortCacheKeyHash` and `api/cache/redis.go:177` `shortNamespaceHash`.
**What:** Byte-identical functions (sha256 → hex → first 12 chars) in the same `cache` package under two names. A reader must diff them to confirm they're the same; a future change to truncation length or algorithm must be made in both places.
**Now:** both confirmed present.
**Fix:** delete one; have the survivor (renamed `shortHash`) serve both call sites.

## #28 — `envInt` defined three times across packages  · P3
**Where:** `api/cache/lru.go:183` `envInt`, `api/middleware/ratelimit.go:188` `envInt`, `api/main.go:54` `getEnvInt`.
**What:** Three near-identical env-int parsers. Two now use `strconv.Atoi` (ratelimit rejects `<=0` + trims; lru rejects `<0`); `main.go`'s `getEnvInt` is the third copy. Minor behavioral divergence (trim/<=0 handling) that could bite during a refactor. (Note: `getEnvInt` itself was hardened against overflow in #6 — this item is about the dedup, not the overflow.)
**Now:** all three copies confirmed present.
**Fix:** move one canonical `envInt` (the ratelimit variant — strictest, with trim) into a small internal package; replace all three.

## #29 — `Embedder` and `WarmupModel` are the same interface under two names  · P3
**Where:** `api/handlers/embeddings.go:24` `Embedder` and `api/lifecycle/warmup.go:11` `WarmupModel`.
**What:** Byte-identical single-method interfaces — `Embed(ctx, []string) ([][]float32, error)`. One implementor (`*embed.TEIClient`); the `handlers.Embedder(teiClient)` conversion in `main.go` is a no-op because `TEIClient` satisfies both. Two names for one contract forces readers to verify equivalence and new consumers to pick which to depend on.
**Now:** both confirmed present.
**Fix:** define one interface (e.g. `embed.Embedder` where `TEIClient` lives); have both `handlers` and `lifecycle` import it; drop the two duplicate declarations.

## #30 — `lifecycle.defaultState` singleton is redundant  · P3
**Where:** `api/lifecycle/state.go:34,43` — `var defaultState = NewState()` + `DefaultState()`.
**What:** The package exposes both `NewState()` (fresh) and `DefaultState()` (singleton). `main.go` is the only production caller and uses `DefaultState()`; `NewState()` is used only by tests and the package-var initializer. Two construction paths with no reason for the split invite the classic bug where one caller mutates the singleton while another holds an independent `NewState()` instance. The singleton also hides the construction site, muddying startup-order reasoning.
**Now:** confirmed present.
**Fix:** drop the singleton — remove `defaultState`/`DefaultState()`; have `main.go` call `lifecycle.NewState()` explicitly and pass it via constructor injection (it already does). Tests already use `NewState()`.

## #31 — `openapi.m()` silently swallows non-string keys  · P3
**Where:** `api/openapi/spec.go:137,142` — `func m(pairs ...any)` does `key, ok := pairs[i].(string); if !ok { continue }`.
**What:** The hand-rolled OpenAPI builder's core helper silently skips a key-value pair on a type error. Because every caller builds map literals from nested `m(...)` calls, a typo or misordered argument in the ~120 lines of nested `paths()`/`components()` calls drops a schema field with no error and no test to catch it. Go's own map literals would catch this at compile time; the variadic-`any` escape hatch trades static type checking for inline convenience in a generator whose entire output is type-bearing data.
**Now:** confirmed present (`continue` still on line 142).
**Fix:** either replace `m()` with a typed Schema builder (compile-time key checking), or change `continue` to a `panic`/logged warning so a misordered argument fails loudly; add a test that decodes `Spec()` and asserts the expected property names per schema.

---

## Scope note

Not in this issue: **#7** (`/metrics` public — likely accepted for scrape), **#22** (OpenAPI omits `/embeddings/batch` — separate spec-drift item), **#16** (`/v1/traces` scoping — mitigated by fail-closed auth, single-tenant by design). All eight items above are P3 housekeeping; safe to batch into one cleanup PR. This issue documents the problems only — implementation to be proposed/PR'd separately.

