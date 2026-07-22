#!/bin/bash
sed -i 's/_, totalControl, durationsControl := runCorpusBurst/totalControl, durationsControl := runCorpusBurst/' embed/coalesce_baseline_test.go
sed -i 's/_, totalCo, durationsCo := runCorpusBurst/totalCo, durationsCo := runCorpusBurst/' embed/coalesce_baseline_test.go

sed -i 's/setupQueueTestServer(t, 2, 32)/setupQueueTestServer(t, 2)/' fd_v2_queue_integration_test.go
sed -i 's/setupQueueTestServer(t, 5, 32)/setupQueueTestServer(t, 5)/' fd_v2_queue_integration_test.go
# Just in case, let's fix the other two instances
sed -i 's/setupQueueTestServer(t, 10, 32)/setupQueueTestServer(t, 10)/' fd_v2_queue_integration_test.go
