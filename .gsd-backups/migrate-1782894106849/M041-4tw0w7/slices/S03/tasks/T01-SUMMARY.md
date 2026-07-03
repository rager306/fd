---
id: T01
parent: S03
milestone: M041-4tw0w7
key_files:
  - api/buildinfo/info.go
  - api/buildinfo/info_test.go
  - api/main.go
  - api/Dockerfile
  - Dockerfile.onnx
  - benchmark-results/m041-s03-t01-go-test.txt
  - benchmark-results/m041-s03-t01-lint.txt
  - benchmark-results/m041-s03-t01-govulncheck.txt
  - benchmark-results/m041-s03-t01-ldflags-build.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T06:24:50.395Z
blocker_discovered: false
---

# T01: Added buildinfo metadata package and ldflags wiring for version, build hash, and build date.

**Added buildinfo metadata package and ldflags wiring for version, build hash, and build date.**

## What Happened

Created `api/buildinfo` with `Info`, stable defaults (`fd-api`, `dev`, `unknown`), and `Uptime()` based on `StartedAt`. Added tests for defaulting, preserving provided metadata, increasing uptime, and zero-value uptime. Added main package ldflags variables `Version`, `BuildHash`, and `BuildDate` using buildinfo defaults so upcoming `/version` and `/info` endpoints have a stable source. Updated `api/Dockerfile` and `Dockerfile.onnx` with `VERSION`, `BUILD_HASH`, and `BUILD_DATE` build args and `-X main.Version`, `-X main.BuildHash`, `-X main.BuildDate` ldflags. Verified a local `go build` with those ldflags succeeds.

## Verification

Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. `go build -trimpath -ldflags='-X main.Version=2.0.0 -X main.BuildHash=abc1234 -X main.BuildDate=2026-06-13T00:00:00Z'` exits 0. Evidence files: benchmark-results/m041-s03-t01-go-test.txt, benchmark-results/m041-s03-t01-lint.txt, benchmark-results/m041-s03-t01-govulncheck.txt, benchmark-results/m041-s03-t01-ldflags-build.txt. GitNexus detect_changes reports LOW risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok including buildinfo | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go build -trimpath -ldflags='-X main.Version=2.0.0 -X main.BuildHash=abc1234 -X main.BuildDate=2026-06-13T00:00:00Z' -o /tmp/fd-api-ldflags-test .` | 0 | ✅ pass: ldflags build succeeds | 180000ms |

## Deviations

The plan's verification path `go test ./api/buildinfo/...` is root-relative and not valid for this repo's Go module layout; equivalent executable verification is `cd api && go test ./buildinfo`. The plan named `Dockerfile`; the active API image file is `api/Dockerfile`, and the opt-in root `Dockerfile.onnx` was updated for consistency.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T01 changes are not committed yet. The built test binary `/tmp/fd-api-ldflags-test` was left in `/tmp` because the harness blocked recursive cleanup traps; it is outside the repository.

## Files Created/Modified

- `api/buildinfo/info.go`
- `api/buildinfo/info_test.go`
- `api/main.go`
- `api/Dockerfile`
- `Dockerfile.onnx`
- `benchmark-results/m041-s03-t01-go-test.txt`
- `benchmark-results/m041-s03-t01-lint.txt`
- `benchmark-results/m041-s03-t01-govulncheck.txt`
- `benchmark-results/m041-s03-t01-ldflags-build.txt`
