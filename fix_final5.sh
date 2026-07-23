#!/bin/bash
set -e
# Restore gosec and ineffassign correctly without messing up return types
sed -i 's/data, err := os.ReadFile(root)/data, err := os.ReadFile(filepath.Clean(root)) \/\/nolint:gosec \/\/ G304: loading local test fixture/' api/embed/coalesce_baseline_test.go
sed -i 's/rng := rand.New(rand.NewSource(42))/rng := rand.New(rand.NewSource(42)) \/\/nolint:gosec \/\/ G404: test data shuffling/' api/embed/coalesce_baseline_test.go
sed -i '/calls = 0/d' api/embed/coalesce_baseline_test.go
sed -i '/totalTexts = 0/d' api/embed/coalesce_baseline_test.go

# Nolint for unparam
sed -i 's/func runCorpusBurst/\/\/nolint:unparam \/\/ calls is used for tests\nfunc runCorpusBurst/' api/embed/coalesce_baseline_test.go

sed -i 's/func setupQueueTestServer/\/\/nolint:unparam \/\/ batchSize varies in tests\nfunc setupQueueTestServer/' api/fd_v2_queue_integration_test.go

sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/' api/fd_v2_queue_integration_test.go
sed -i 's/(t \*testing.T, queueCap int, batchSize int)/(t *testing.T, queueCap, batchSize int)/' api/fd_v2_queue_integration_test.go
sed -i 's/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, nil)/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, http.NoBody)/g' api/fd_v2_queue_integration_test.go
sed -i 's/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", nil)/httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", http.NoBody)/g' api/fd_v2_queue_integration_test.go
sed -i 's/func postQueue(t \*testing.T, r http.Handler, body string) \*httptest.ResponseRecorder {/func postQueue(_ \*testing.T, r http.Handler, body string) \*httptest.ResponseRecorder {/' api/fd_v2_queue_integration_test.go
