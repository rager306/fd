1. Add explicit error handling/ignores to satisfy `errcheck`
   - In `api/queue/worker_test.go`, wrap `defer store.Close()` in a func that ignores the error: `defer func() { _ = store.Close() }()`. Do this for all 5 occurrences lines 63, 95, 132, 165, 208.
   - In `api/main_test.go`, lines 277, 289, wrap `defer c.Close()` in `defer func() { _ = c.Close() }()`.
   - In `api/main.go`, line 459, wrap `defer resultStore.Close()` in `defer func() { _ = resultStore.Close() }()`.
   - In `api/fd_v2_queue_integration_test.go`, line 37, wrap `t.Cleanup(func() { store.Close() })` to `t.Cleanup(func() { _ = store.Close() })`.
   - In `api/embed/coalescedembedder_test.go`, line 86, wrap `co.Embed(context.Background(), []string{"x"})` to `_, _ = co.Embed(context.Background(), []string{"x"})` or `_ = co.Embed...` if it returns 1 value. (Need to check return type).

2. Fix `goconst` issues
   - Extract `"list"` string to a constant in `api/handlers/embeddings.go` (line 263), `api/handlers/embeddings_integration_test.go` (line 179), and `api/handlers/queue_handlers.go` (line 127). Wait, I can just append `//nolint:goconst // constant not needed` or add `const objectList = "list"`. It's better to add the nolint. Wait, we need an explanation. `//nolint:goconst // API response literal`.
   - Same for `"tier"` in `api/observability/metrics.go` line 88.

3. Fix `gocritic` issues
   - `paramTypeCombine`: In `api/embed/coalesce_baseline_test.go:70` and `api/fd_v2_queue_integration_test.go:32`
   - `ifElseChain`: In `api/embed/coalescedembedder.go:141`
   - `httpNoBody`: In `api/fd_v2_queue_integration_test.go` lines 81 and 142.
   - `exitAfterDefer`: In `api/main.go:523` -> add `//nolint:gocritic // exitAfterDefer: intentional fast crash` and manually invoke `recoveryCancel()` before `os.Exit(1)`.

4. Fix `gocyclo` issue
   - In `api/main.go:231`, add `//nolint:gocyclo // main function setup block`

5. Fix `gosec` issues
   - `api/embed/coalesce_baseline_test.go:42`: `data, err := os.ReadFile(root)` -> use `filepath.Clean(root)` or add `//nolint:gosec // G304: loading local test fixture`
   - `api/embed/coalesce_baseline_test.go:86`: `rng := rand.New(rand.NewSource(42))` -> use `crypto/rand` or we can just ignore `//nolint:gosec // G404: test data generation doesn't need crypto/rand`. Wait, memory says "Avoid using math/rand even in test or benchmark files; use fixed test data or crypto/rand to prevent linting failures." So I'll rewrite to use a fast non-crypto rand or crypto/rand.

6. Fix `govet` issue
   - `api/cache/tiered.go:233`: `defer tc.recordLookupDuration(time.Since(started))` -> `defer func() { tc.recordLookupDuration(time.Since(started)) }()`

7. Fix `ineffassign`
   - `api/embed/coalesce_baseline_test.go` lines 72, 74. Remove assignments.

8. Fix `revive` issues
   - `api/cache/tiered.go:29`: rename `CacheObserver` to `Observer`.
   - `api/embed/tei.go:251`: rename `ObserveBatchFill` to `RecordBatchFill` or something else? Wait, if it's exported, I'll need to check where it's used. Let's find out.
   - `api/observability/metrics.go:223`: add a comment.
   - `api/queue/types.go:21`: add a comment.

9. Fix `staticcheck`
   - `api/queue/worker.go:164`: `indexByID = append(indexByID, &batch[i])` is never used.

10. Fix `unparam`
    - `api/embed/coalescedembedder_test.go:42`: `i is unused` -> change to `_ int` or `_`.
    - `api/fd_v2_queue_integration_test.go:57`: `t *testing.T` is unused -> change `t` to `_` or remove.

11. Fix `unused`
    - `api/observability/metrics.go:56`: `redisSizeTimeout` -> remove.
