sed -i 's/func postQueue(t \*testing.T, r http.Handler, body string) \*httptest.ResponseRecorder {/func postQueue(_ \*testing.T, r http.Handler, body string) \*httptest.ResponseRecorder {/g' api/fd_v2_queue_integration_test.go
sed -i 's/go func(i int) {/go func(_ int) {/g' api/embed/coalescedembedder_test.go
