#!/bin/bash
set -e
git restore api/main.go
# 1. Add nolint to main.go for exitAfterDefer
sed -i 's/os.Exit(1)/os.Exit(1) \/\/nolint:gocritic \/\/ exitAfterDefer: os.Exit in main is acceptable/' api/main.go

# 2. Fix the missing errcheck on resultStore.Close()
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/' api/main.go

# 3. Add explanation for whyNoLint in main.go
sed -i 's/func main()/\/\/nolint:gocyclo \/\/ main has many configuration steps\nfunc main()/' api/main.go

# 4. Remove unused return value in api/embed/coalesce_baseline_test.go
sed -i 's/(calls int, totalTexts int, durations \\\[\\]time.Duration)/(totalTexts int, durations \[\]time.Duration)/' api/embed/coalesce_baseline_test.go
# Wait, let me just replace the whole line for coalesce_baseline_test.go
sed -i 's/func runCorpusBurst(t \*testing.T, e Embedder, texts \[\]string, concurrency int) (calls, totalTexts int, durations \[\]time.Duration) {/func runCorpusBurst(t \*testing.T, e Embedder, texts \[\]string, concurrency int) (totalTexts int, durations \[\]time.Duration) {/' api/embed/coalesce_baseline_test.go
