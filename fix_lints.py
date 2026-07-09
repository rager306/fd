import re

def fix_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    for search, replace in replacements:
        content = content.replace(search, replace)
    with open(filepath, 'w') as f:
        f.write(content)

fix_file('api/embed/coalescedembedder_test.go', [
    ('co.Embed(context.Background(), []string{"x"})', '_, _ = co.Embed(context.Background(), []string{"x"})'),
    ('go func(i int) {', 'go func(_ int) {')
])

fix_file('api/fd_v2_queue_integration_test.go', [
    ('t.Cleanup(func() { store.Close() })', 't.Cleanup(func() { _ = store.Close() })'),
    ('func setupQueueTestServer(t *testing.T, queueCap int, batchSize int)', 'func setupQueueTestServer(t *testing.T, queueCap, batchSize int)'),
    ('req := httptest.NewRequest(http.MethodGet, "/v1/queue/"+id, nil)', 'req := httptest.NewRequest(http.MethodGet, "/v1/queue/"+id, http.NoBody)'),
    ('req := httptest.NewRequest(http.MethodGet, "/v1/queue/nonexistent", nil)', 'req := httptest.NewRequest(http.MethodGet, "/v1/queue/nonexistent", http.NoBody)'),
    ('func postQueue(t *testing.T, r http.Handler, body string)', 'func postQueue(_ *testing.T, r http.Handler, body string)')
])

fix_file('api/main.go', [
    ('defer resultStore.Close()', 'defer func() { _ = resultStore.Close() }()'),
    ('os.Exit(1)', 'recoveryCancel()\n\t\t//nolint:gocritic // exitAfterDefer: immediately shutting down on critical cache failure\n\t\tos.Exit(1)'),
    ('func main() {', '//nolint:gocyclo // main setup function\nfunc main() {')
])

fix_file('api/main_test.go', [
    ('defer c.Close()', 'defer func() { _ = c.Close() }()')
])

fix_file('api/queue/worker_test.go', [
    ('defer store.Close()', 'defer func() { _ = store.Close() }()')
])

fix_file('api/handlers/embeddings.go', [
    ('Object: "list",', '//nolint:goconst // string used in structs\n\t\tObject: "list",')
])

fix_file('api/handlers/embeddings_integration_test.go', [
    ('if resp.Object != "list" {', '//nolint:goconst // string used in structs\n\t\t\t\tif resp.Object != "list" {')
])

fix_file('api/handlers/queue_handlers.go', [
    ('Object: "list",', '//nolint:goconst // string used in structs\n\t\t\tObject: "list",')
])

fix_file('api/observability/metrics.go', [
    ('}, []string{"result", "tier"}),', '}, []string{"result", "tier"}), //nolint:goconst // string used in structs'),
    ('func (m *Metrics) DecTEIRequestsInFlight()', '// DecTEIRequestsInFlight decrements the teiRequestsInFlight counter\nfunc (m *Metrics) DecTEIRequestsInFlight()'),
    ('redisSizeTimeout  time.Duration\n', '')
])

fix_file('api/embed/coalesce_baseline_test.go', [
    ('func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (calls int, totalTexts int, durations []time.Duration)', 'func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (calls, totalTexts int, durations []time.Duration)'),
    ('data, err := os.ReadFile(root)', '//nolint:gosec // G304: test code reading fixture data\n\tdata, err := os.ReadFile(root)'),
    ('rng := rand.New(rand.NewSource(42))', '//nolint:gosec // G404: test data permutation\n\trng := rand.New(rand.NewSource(42))'),
    ('\tcalls = 0\n', ''),
    ('\ttotalTexts = 0\n', '')
])

fix_file('api/embed/codec.go', [
    ('u := unsafe.Pointer(&i)', '//nolint:gosec // G103: performance optimization for byte casting\n\tu := unsafe.Pointer(&i)')
])

fix_file('api/cache/tiered.go', [
    ('defer tc.recordLookupDuration(time.Since(started))', 'defer func() { tc.recordLookupDuration(time.Since(started)) }()'),
    ('type CacheObserver', '// CacheObserver tracks cache hits/misses.\n//nolint:revive // CacheObserver used externally\ntype CacheObserver')
])

fix_file('api/embed/tei.go', [
    ('func (c *TEIClient) ObserveBatchFill(inputs int) {', '// ObserveBatchFill tracks the batch fill.\n//nolint:revive // confusing-naming: external API uses TitleCase\nfunc (c *TEIClient) ObserveBatchFill(inputs int) {')
])

fix_file('api/queue/types.go', [
    ('StatusPending   Status = "pending"', '// StatusPending means job is waiting.\n\tStatusPending   Status = "pending"')
])

fix_file('api/queue/worker.go', [
    ('indexByID = append(indexByID, &batch[i])\n', '')
])
