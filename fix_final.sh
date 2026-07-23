#!/bin/bash
set -e

# Fix whyNoLint in api/main.go
# Sometimes gocritic just expects it in the same line or standard format. Wait, it is:
# //nolint:gocyclo // main has many configuration steps
# Maybe we didn't add it correctly, or it's not the right syntax. Let's just drop the nolint and increase gocyclo max in .golangci.yml if possible, or refactor main.
# Better to just increase `min-complexity: 15` to `18` or `20` in .golangci.yml or just use the right nolint:
# Let's change the comment to just `//nolint:gocyclo // main setup` and check again, wait, maybe gocritic's `whyNoLint` requires `//nolint:gocyclo // ...` and earlier it was `//nolint:gocyclo // main function setup` which gave the same error.
# Ah, I replaced `//nolint:gocyclo // main function setup` with `//nolint:gocyclo // main has many configuration steps` in the script, maybe it didn't apply because the file was restored.
sed -i 's/\/\/nolint:gocyclo \/\/ main has many configuration steps/\/\/nolint:gocyclo \/\/ main/' api/main.go # wait, maybe the exact string was changed.

# Fix gosec G118
sed -i 's/recoveryCtx, _ := context.WithCancel(context.Background())/recoveryCtx, recoveryCancel := context.WithCancel(context.Background())\n\tdefer recoveryCancel()/' api/main.go

# Fix unparam in runCorpusBurst
sed -i 's/(calls, totalTexts int, durations \[\\]time.Duration)/(totalTexts int, durations \[\\]time.Duration)/' api/embed/coalesce_baseline_test.go
# wait, there are uses of `calls`?
