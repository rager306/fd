# api/cache/tiered.go
sed -i 's/defer tc.recordLookupDuration(time.Since(started))/defer func() { tc.recordLookupDuration(time.Since(started)) }()/' api/cache/tiered.go
sed -i 's/type CacheObserver/type Observer/' api/cache/tiered.go
sed -i 's/observer CacheObserver/observer Observer/' api/cache/tiered.go
sed -i 's/cache.CacheObserver/cache.Observer/g' api/main.go
sed -i 's/cache.CacheObserver/cache.Observer/g' api/main_test.go

# api/embed/coalescedembedder.go
sed -i 's/if err != nil {/if err != nil {/' api/embed/coalescedembedder.go # We'll need a git replace for switch

# api/embed/coalesce_baseline_test.go
sed -i 's/(calls int, totalTexts int, durations/(calls, totalTexts int, durations/g' api/embed/coalesce_baseline_test.go
sed -i 's/calls = 0/var calls int/g' api/embed/coalesce_baseline_test.go
sed -i 's/totalTexts = 0/var totalTexts int/g' api/embed/coalesce_baseline_test.go
sed -i 's/"math\/rand"/"math\/rand\/v2"/g' api/embed/coalesce_baseline_test.go
sed -i 's/rand.New(rand.NewSource(42))/rand.New(rand.NewPCG(42, 42))/g' api/embed/coalesce_baseline_test.go

# api/fd_v2_queue_integration_test.go
sed -i 's/queueCap int, batchSize int/queueCap, batchSize int/g' api/fd_v2_queue_integration_test.go
sed -i 's/nil)/http.NoBody)/g' api/fd_v2_queue_integration_test.go

# api/embed/tei.go
sed -i 's/ObserveBatchFill(/ObserveBatchFillExported(/g' api/embed/tei.go

# api/observability/metrics.go
sed -i 's/func (m \*Metrics) DecTEIRequestsInFlight()/\/\/ DecTEIRequestsInFlight decrements the in-flight TEI requests counter.\nfunc (m \*Metrics) DecTEIRequestsInFlight()/' api/observability/metrics.go
sed -i '/redisSizeTimeout time.Duration/d' api/observability/metrics.go

# api/queue/types.go
sed -i 's/StatusPending   Status = "pending"/\/\/ StatusPending indicates the item is waiting to be processed.\n	StatusPending   Status = "pending"/' api/queue/types.go

# api/queue/worker.go
sed -i '/indexByID = append(indexByID, &batch\[i\])/d' api/queue/worker.go
