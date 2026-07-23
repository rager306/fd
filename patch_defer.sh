#!/bin/bash
set -e

# Replace os.Exit(1) with return because this is the main function and we can just let it exit normally or panic, but better to just use return and set exit status or just return. Actually, if we want non-zero exit, we can't use `os.Exit` directly in `main` if we have defers we want to run.
# Or we can just drop `defer recoveryCancel()` and call it manually. Wait, `recoveryCancel()` is small. Let's just remove `os.Exit(1)` and change main to a function `run() error` but that's too big.
# Instead of `os.Exit(1)`, we can `recoveryCancel()` explicitly before exiting.
sed -i 's/os.Exit(1)/recoveryCancel()\n\t\tos.Exit(1)/g' api/main.go
