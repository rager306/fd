#!/bin/bash
set -e
cd api

# 1. Error return value of `co.Embed` is not checked (errcheck)
sed -i 's/co.Embed(context.Background(), \[\]string{"x"})/_, _ = co.Embed(context.Background(), []string{"x"})/g' embed/coalescedembedder_test.go

# 2. Error return value of `store.Close` is not checked (errcheck)
sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/g' fd_v2_queue_integration_test.go

sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/g' main.go

sed -i 's/defer c.Close()/defer func() { _ = c.Close() }()/g' main_test.go

sed -i 's/defer store.Close()/defer func() { _ = store.Close() }()/g' queue/worker_test.go
