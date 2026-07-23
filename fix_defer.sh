#!/bin/bash
set -e
# Revert patch_defer.sh
git checkout api/main.go

sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/' api/main.go
sed -i 's/func main()/ \/\/nolint:gocyclo\nfunc main()/' api/main.go
