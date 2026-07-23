#!/bin/bash
sed -i '/\/\/ StatusPending, StatusCompleted, StatusFailed represent queue states/d' api/queue/types.go
sed -i 's/\/\/ Status constants for queue items.//' api/queue/types.go
sed -i 's/\/\/ Status is the lifecycle of a submitted queue item./\/\/ Status is the lifecycle of a submitted queue item.\n\/\/ StatusPending is the initial state./' api/queue/types.go
