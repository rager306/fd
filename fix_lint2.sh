sed -i 's/func runCorpusBurst(t \*testing.T, e Embedder, texts \[\]string, concurrency int) (calls int, totalTexts int, durations \[\]time.Duration)/func runCorpusBurst(t \*testing.T, e Embedder, texts \[\]string, concurrency int) (calls, totalTexts int, durations \[\]time.Duration)/g' api/embed/coalesce_baseline_test.go
sed -i 's/func setupQueueTestServer(t \*testing.T, queueCap int, batchSize int)/func setupQueueTestServer(t \*testing.T, queueCap, batchSize int)/g' api/fd_v2_queue_integration_test.go
sed -i 's/req := httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, nil)/req := httptest.NewRequest(http.MethodGet, "\/v1\/queue\/"+id, http.NoBody)/g' api/fd_v2_queue_integration_test.go
sed -i 's/req := httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", nil)/req := httptest.NewRequest(http.MethodGet, "\/v1\/queue\/nonexistent", http.NoBody)/g' api/fd_v2_queue_integration_test.go
sed -i 's/os.Exit(1)/recoveryCancel(); os.Exit(1) \/\/nolint:gocritic \/\/ exitAfterDefer: need to cancel contexts manually/g' api/main.go
sed -i 's/defer tc.recordLookupDuration(time.Since(started))/defer func() { tc.recordLookupDuration(time.Since(started)) }()/g' api/cache/tiered.go
sed -i 's/calls = 0/ /g' api/embed/coalesce_baseline_test.go
sed -i 's/totalTexts = 0/ /g' api/embed/coalesce_baseline_test.go
sed -i 's/type CacheObserver func/type Observer func/g' api/cache/tiered.go
sed -i 's/\/\/ ObserveBatchFill is/\/\/ ObserveFill is/g' api/embed/tei.go
sed -i 's/func (c \*TEIClient) ObserveBatchFill(inputs int)/func (c \*TEIClient) ObserveFill(inputs int)/g' api/embed/tei.go
sed -i 's/\/\/ IncTEIRequestsInFlight and DecTEIRequestsInFlight/\/\/ IncTEIRequestsInFlight and DecTEIRequestsInFlight/g' api/observability/metrics.go
