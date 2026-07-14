#!/bin/bash
set -e

echo "Fixing api/queue/worker_test.go"
sed -i 's/defer store.Close()/defer func() { _ = store.Close() }()/g' api/queue/worker_test.go

echo "Fixing api/main_test.go"
sed -i 's/defer c.Close()/defer func() { _ = c.Close() }()/g' api/main_test.go

echo "Fixing api/main.go errcheck"
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/g' api/main.go

echo "Fixing api/fd_v2_queue_integration_test.go errcheck"
sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/g' api/fd_v2_queue_integration_test.go

echo "Fixing api/embed/coalescedembedder_test.go errcheck"
sed -i 's/co.Embed(context.Background(), \[\]string{"x"})/_, _ = co.Embed(context.Background(), []string{"x"})/g' api/embed/coalescedembedder_test.go

echo "Fixing api/fd_v2_queue_integration_test.go paramTypeCombine"
sed -i 's/func setupQueueTestServer(t \*testing.T, queueCap int, batchSize int)/func setupQueueTestServer(t \*testing.T, queueCap, batchSize int)/g' api/fd_v2_queue_integration_test.go

echo "Fixing api/fd_v2_queue_integration_test.go httpNoBody"
sed -i 's/, nil)/, http.NoBody)/g' api/fd_v2_queue_integration_test.go

echo "Fixing api/embed/coalesce_baseline_test.go paramTypeCombine"
sed -i 's/calls int, totalTexts int, durations \[\]time.Duration/calls, totalTexts int, durations []time.Duration/g' api/embed/coalesce_baseline_test.go
