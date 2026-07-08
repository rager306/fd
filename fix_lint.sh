#!/bin/bash

# fix errcheck in queue/worker_test.go
sed -i 's/defer store.Close()/defer func() { _ = store.Close() }()/' api/queue/worker_test.go

# fix errcheck in main_test.go
sed -i 's/defer c.Close()/defer func() { _ = c.Close() }()/' api/main_test.go

# fix errcheck in main.go
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/' api/main.go

# fix errcheck in fd_v2_queue_integration_test.go
sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/' api/fd_v2_queue_integration_test.go

# fix errcheck in coalescedembedder_test.go
sed -i 's/co.Embed(context.Background(), \[\]string{"x"})/_, _ = co.Embed(context.Background(), []string{"x"})/' api/embed/coalescedembedder_test.go
