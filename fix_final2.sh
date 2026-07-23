#!/bin/bash
set -e
# Restore main.go, fix defer issue properly.
# The error was: api/main.go:507:38: G118: context cancellation function returned by WithCancel/WithTimeout/WithDeadline is not called (gosec)
# This happens because we removed `defer recoveryCancel()` since it was panicking? No, we removed it because there was a compile error because it was used in defer, and also returned by WithCancel. Wait!
# The compile error earlier: `undefined: recoveryCancel`.
# Looking at the original api/main.go around line 507:
# 	recoveryCtx, recoveryCancel := context.WithCancel(context.Background())
# 	defer recoveryCancel()
# But wait, lines 249, 268, 275, 286 also used `recoveryCancel()`.
