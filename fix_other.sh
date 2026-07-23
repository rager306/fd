#!/bin/bash
set -e
# CacheObserver revive
sed -i 's/\/\/ CacheObserver is invoked/\/\/ Observer is invoked/' api/cache/tiered.go

# StatusPending revive
sed -i 's/\/\/ Status constants for queue items./\/\/ StatusPending is the initial state of a queue item./' api/queue/types.go
# Wait, let's just make it a block comment:
# 	// Status constants:
# 	StatusPending   Status = "pending"
# 	StatusCompleted Status = "completed"
sed -i 's/\/\/ StatusPending is the initial state of a queue item./\/\/ Status constants for queue items/' api/queue/types.go # Wait, we changed it earlier to // Status constants for queue items.
sed -i 's/\/\/ Status constants for queue items/\/\/ StatusPending, StatusCompleted, StatusFailed represent queue states/' api/queue/types.go
