#!/bin/bash
set -e

# paramTypeCombine
sed -i 's/(calls int, totalTexts int, durations \\\[\\]time.Duration)/(calls, totalTexts int, durations \[\]time.Duration)/' api/embed/coalesce_baseline_test.go
sed -i 's/(t \*testing.T, queueCap int, batchSize int)/(t *testing.T, queueCap, batchSize int)/' api/fd_v2_queue_integration_test.go

# whyNoLint
sed -i 's/\/\/nolint:gocyclo \/\/ main has many configuration steps/\/\/nolint:gocyclo \/\/ main function setup requires many branches/' api/main.go

# gosec G118
sed -i 's/recoveryCtx, _ := context.WithCancel(context.Background())/recoveryCtx, recoveryCancel := context.WithCancel(context.Background())\n\tdefer recoveryCancel()/' api/main.go
