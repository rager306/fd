#!/bin/bash
set -e
sed -i 's/defer resultStore.Close()/defer func() { _ = resultStore.Close() }()/' api/main.go
sed -i 's/func main()/\/\/nolint:gocyclo\nfunc main()/' api/main.go
