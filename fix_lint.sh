#!/bin/bash
sed -i '523s/recoveryCancel(); os.Exit(1) \/\/nolint:gocritic \/\/ exitAfterDefer: manual cancel/recoveryCancel(); os.Exit(1) \/\/nolint:gocritic \/\/ exitAfterDefer: manual cancel/g' api/main.go
