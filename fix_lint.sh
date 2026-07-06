sed -i 's/co.Embed(context.Background(), \[\]string{"x"})/_, _ = co.Embed(context.Background(), \[\]string{"x"})/g' api/embed/coalescedembedder_test.go
sed -i 's/t.Cleanup(func() { store.Close() })/t.Cleanup(func() { _ = store.Close() })/g' api/fd_v2_queue_integration_test.go
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/g' api/main.go
sed -i 's/defer c.Close()/defer func() { _ = c.Close() }()/g' api/main_test.go
sed -i 's/defer store.Close()/defer func() { _ = store.Close() }()/g' api/queue/worker_test.go

sed -i 's/Object: "list"/Object: constList/g' api/handlers/embeddings.go
sed -i 's/if resp.Object != "list"/if resp.Object != constList/g' api/handlers/embeddings_integration_test.go
sed -i 's/Object: "list"/Object: constList/g' api/handlers/queue_handlers.go
sed -i 's/\[\]string{"result", "tier"}/\[\]string{"result", constTier}/g' api/observability/metrics.go
