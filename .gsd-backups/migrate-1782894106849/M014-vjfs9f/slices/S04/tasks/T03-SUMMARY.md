---
id: T03
parent: S04
milestone: M014-vjfs9f
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
  - benchmark-results/fd-benchmark-m014-comparison.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T04:35:53.416Z
blocker_discovered: false
---

# T03: Final M014 verification gates passed after the last code change.

**Final M014 verification gates passed after the last code change.**

## What Happened

Ran the final verification gates after the last code change. Default Go tests passed, pinned GolangCI-Lint reported 0 issues, tagged `hf_tokenizers` tests passed, benchmark.py compiles, artifact hygiene checks passed with zero raw fixed-probe text leaks, no native/ONNX binaries are tracked, default API health is ok, the tagged ONNX port is stopped, Docker Compose services are healthy, no bg_shell processes remain, and GitNexus scope is low/non-code for pending S04 artifacts.

## Verification

Fresh final gates passed: Go tests 78 passed, lint 0 issues, tagged tests 20 passed, artifact hygiene/leak checks pass, tracked native binaries 0, default health ok, ONNX server stopped, no background processes, GitNexus low artifact-only scope.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — 78 passed in 4 packages | 7000ms |
| 2 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 7000ms |
| 3 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 packages | 6900ms |
| 4 | `python3 -m py_compile benchmark.py && artifact hygiene/leak/tracked binary checks` | 0 | ✅ pass — artifact_hygiene=pass; raw_probe_text_leaks=0; tracked_native_binaries=0 | 0ms |
| 5 | `curl health/default + curl port 18000 cleanup + docker compose ps` | 0 | ✅ pass — default API ok; ONNX server stopped; compose services healthy | 0ms |
| 6 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |
| 7 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low scope, no changed symbols | 0ms |

## Deviations

None.

## Known Issues

GitNexus detect currently shows low artifact-only scope before final S04 commit. TEI Compose API remains healthy; tagged ONNX server is stopped.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
- `benchmark-results/fd-benchmark-m014-comparison.txt`
