#!/bin/bash
set -e

# Fix errcheck in api/embed/coalescedembedder_test.go:86
sed -i 's/co.Embed(context.Background(), \[\]string{"x"})/_, _ = co.Embed(context.Background(), []string{"x"})/g' api/embed/coalescedembedder_test.go

# Fix errcheck in api/fd_v2_queue_integration_test.go:37
sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/g' api/fd_v2_queue_integration_test.go

# Fix errcheck in api/main.go:459
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/g' api/main.go

# Fix errcheck in api/main_test.go:277 and 289
sed -i 's/defer c.Close()/defer func() { _ = c.Close() }()/g' api/main_test.go

# Fix errcheck in api/queue/worker_test.go:63, 95, 132, 165, 208
sed -i 's/defer store.Close()/defer func() { _ = store.Close() }()/g' api/queue/worker_test.go

# Fix goconst in api/handlers/embeddings.go:263
sed -i 's/Object: "list",/Object: objectList,/g' api/handlers/embeddings.go
sed -i '/const (/a \\tobjectList = "list"' api/handlers/embeddings.go

# Fix goconst in api/handlers/embeddings_integration_test.go:179
sed -i 's/if resp.Object != "list" {/if resp.Object != objectList {/g' api/handlers/embeddings_integration_test.go

# Fix goconst in api/handlers/queue_handlers.go:127
sed -i 's/Object: "list",/Object: objectList,/g' api/handlers/queue_handlers.go
sed -i '/const (/a \\tobjectList = "list"' api/handlers/queue_handlers.go

# Fix goconst in api/observability/metrics.go:88
sed -i 's/\[\]string{"result", "tier"}/[]string{"result", labelTier}/g' api/observability/metrics.go
sed -i '/const (/a \\tlabelTier = "tier"' api/observability/metrics.go
if ! grep -q 'labelTier' api/observability/metrics.go; then
    sed -i '/import (/a \\nconst labelTier = "tier"' api/observability/metrics.go
fi


# Fix gocritic paramTypeCombine in api/embed/coalesce_baseline_test.go:70
sed -i 's/(calls int, totalTexts int, durations \[\]time.Duration)/(calls, totalTexts int, durations []time.Duration)/g' api/embed/coalesce_baseline_test.go

# Fix gocritic ifElseChain in api/embed/coalescedembedder.go:141
cat << 'INNEREOF' > patch_ifelse.txt
<<<<<<< SEARCH
		if err != nil {
			j.result <- coalescedResult{err: err}
		} else if cursor+n <= len(embs) {
			slice := make([][]float32, n)
			copy(slice, embs[cursor:cursor+n])
			j.result <- coalescedResult{embeddings: slice}
		} else {
			j.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}
		}
=======
		switch {
		case err != nil:
			j.result <- coalescedResult{err: err}
		case cursor+n <= len(embs):
			slice := make([][]float32, n)
			copy(slice, embs[cursor:cursor+n])
			j.result <- coalescedResult{embeddings: slice}
		default:
			j.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}
		}
>>>>>>> REPLACE
INNEREOF
patch -p1 < patch_ifelse.txt || true # use replace diff later if fails

# Fix gocritic paramTypeCombine in api/fd_v2_queue_integration_test.go:32
sed -i 's/(t \*testing.T, queueCap int, batchSize int)/(t *testing.T, queueCap, batchSize int)/g' api/fd_v2_queue_integration_test.go

# Fix gocritic httpNoBody in api/fd_v2_queue_integration_test.go:81 and 142
sed -i 's/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, nil)/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, http.NoBody)/g' api/fd_v2_queue_integration_test.go
sed -i 's/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", nil)/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", http.NoBody)/g' api/fd_v2_queue_integration_test.go

# Fix gocritic exitAfterDefer in api/main.go:523
# Change `defer recoveryCancel()` or similar to be called directly, or structure it differently.
# But looking at line 523 `os.Exit(1)` inside `main` while defers exist. Let's just avoid `os.Exit(1)` and return instead if it's in a func, but in `main` we can call `closeResource` directly.
# wait, actually the linter complains `os.Exit` will exit and defer won't run.
# Let's check `api/main.go:515-530`.
#!/bin/bash
set -e

# Fix exitAfterDefer in api/main.go:523
# Let's review the main.go around line 523
